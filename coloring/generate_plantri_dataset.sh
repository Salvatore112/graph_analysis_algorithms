#!/usr/bin/env bash
set -euo pipefail

PLANTRI_BIN="${PLANTRI:-plantri}"

TRI_NS=${TRI_NS:-"10 12 14"}
TRI_FRACTION=${TRI_FRACTION:-""}
TRI_HEAD=${TRI_HEAD:-1000}


PP_NS=${PP_NS:-"10 12"}
PP_FRACTION=${PP_FRACTION:-"1/5000"}
PP_HEAD=${PP_HEAD:-500}


TIMEOUT_SECS=${TIMEOUT_SECS:-60}

REGEN=${REGEN:-0}


SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
OUT_DIR="${SCRIPT_DIR}/dataset"
mkdir -p "${OUT_DIR}"

if [[ -x "${PLANTRI_BIN}" ]]; then
  PLANTRI="${PLANTRI_BIN}"
elif command -v "${PLANTRI_BIN}" >/dev/null 2>&1; then
  PLANTRI="$(command -v "${PLANTRI_BIN}")"
else
  echo "Не найден plantri. Укажи PLANTRI=/abs/path/to/plantri или добавь его в PATH." >&2
  exit 1
fi

if command -v timeout >/dev/null 2>&1; then
  HAVE_TIMEOUT=1
else
  HAVE_TIMEOUT=0
fi

echo "Использую plantri: ${PLANTRI}"
echo "Каталог вывода:    ${OUT_DIR}"
echo

gen_file() {
  local cmd="$1" outfile="$2" head_limit="${3:-0}"

  if [[ -f "${outfile}" && "${REGEN}" -ne 1 ]]; then
    echo "[skip] ${outfile} уже существует (REGEN=0)"
    return
  fi

  echo "[gen ] ${outfile}"
  local tmp="${outfile}.tmp"

  if [[ "${TIMEOUT_SECS}" -gt 0 && "${HAVE_TIMEOUT}" -eq 1 ]]; then
    if ! timeout "${TIMEOUT_SECS}s" bash -c "${cmd} > '${tmp}'"; then
      echo "[warn] команда превысила тайм-аут ${TIMEOUT_SECS}s: ${cmd}" >&2
      [[ -f "${tmp}" ]] || : > "${tmp}"
    fi
  else
    bash -c "${cmd} > '${tmp}'"
  fi

  if [[ "${head_limit}" -gt 0 ]]; then
    head -n "${head_limit}" "${tmp}" > "${outfile}"
    rm -f "${tmp}"
  else
    mv "${tmp}" "${outfile}"
  fi
}

for n in ${TRI_NS}; do
  out="${OUT_DIR}/tri_${n}.g6"
  if [[ -n "${TRI_FRACTION}" ]]; then
    gen_file "'${PLANTRI}' -g ${n} ${TRI_FRACTION}" "${out}" "${TRI_HEAD}"
  else
    gen_file "'${PLANTRI}' -g ${n}" "${out}" "${TRI_HEAD}"
  fi
done

for n in ${PP_NS}; do
  out="${OUT_DIR}/ppoly_${n}.g6"
  if [[ -n "${PP_FRACTION}" ]]; then
    gen_file "'${PLANTRI}' -pg ${n} ${PP_FRACTION}" "${out}" "${PP_HEAD}"
  else
    gen_file "'${PLANTRI}' -pg ${n}" "${out}" "${PP_HEAD}"
  fi
done

echo
echo "Готово. Файлы *.g6 лежат в ${OUT_DIR}."
echo "Советы: для крупных n обязательно задавай доли TRI_FRACTION/PP_FRACTION или увеличивай TIMEOUT_SECS."
