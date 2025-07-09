#!/bin/bash

echo "🔍 Ejecutando tests..."

# Ejecuta los tests con salida detallada (-v)
go test ./... -v

# Guarda el código de salida
EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
  echo "✅ Todos los tests pasaron correctamente."
else
  echo "❌ Algunos tests fallaron."
fi

exit $EXIT_CODE
