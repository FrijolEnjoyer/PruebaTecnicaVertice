#!/bin/bash

echo "🔍 Ejecutando tests con cobertura..."

go test -v -coverprofile=coverage.out ./...


EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
  echo "✅ Todos los tests pasaron correctamente."
  echo "📊 Reporte de cobertura:"
  go tool cover -func=coverage.out
else
  echo "❌ Algunos tests fallaron."
fi


exit $EXIT_CODE
