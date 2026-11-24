# Actualizar el contrato OpenAPI

> **Estado:** Aprobado

El contrato vive en un solo archivo, [openapi.yaml](/docs/reference/api/openapi.yaml) (OpenAPI 3.2.0). No se genera
desde el código ni se sirve: se edita a mano junto con el cambio.

## Cuándo

- Un endpoint nuevo o un campo nuevo en un cuerpo
- Un código de error nuevo en el mapeo de errores del adaptador HTTP ([handlers.go](/apps/api/adapters/httpapi/handlers.go), función `fail`)
- Un formato o límite nuevo en las validaciones de [service.go](/apps/api/core/service.go) (función `Send`)

## Pasos

1. Rutas en `paths`, con `operationId`, `tags` y una respuesta por código
2. Schemas reutilizables en `components/schemas`; los `pattern` copian las expresiones de `service.go`
3. Errores repetidos en `components/responses`
4. Descripciones en español, sin punto final, como el resto del archivo
5. Si el cambio toca el portal, actualizar también [contrato-portal.md](/docs/reference/api/contrato-portal.md)

## Validar

Sin herramientas en el repo. Para revisar la sintaxis YAML:

```bash
ruby -ryaml -e 'YAML.load_file("docs/reference/api/openapi.yaml")'
```
