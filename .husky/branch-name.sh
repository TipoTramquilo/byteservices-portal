# Validación de nombre de rama, compartida por post-checkout (avisa) y pre-push (bloquea)
# Ramas protegidas libres: main, develop, staging. El resto: <tipo>/<descripción-en-kebab>
BRANCH_PATTERN='^(feature|fix|refactor|hotfix|chore|docs)/[a-z0-9-]+$'

branch_ok() {
  case "$1" in
    main | develop | staging | HEAD) return 0 ;;
  esac
  echo "$1" | grep -qE "$BRANCH_PATTERN"
}

branch_help() {
  echo "  Formato esperado: <tipo>/<descripción>"
  echo "  Tipos válidos: feature|fix|refactor|hotfix|chore|docs"
  echo "  Ejemplo: feature/ci-robusto"
  echo "  Crearla bien: git checkout -b <tipo>/<descripción>"
  echo "  Renombrar la actual: git branch -m <tipo>/<descripción>"
}
