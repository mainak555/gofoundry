# Code Standards

## GoDoc and Comments

- Every exported type, function, method, and interface must have a GoDoc comment.
- Comment text starts with the symbol name and explains behavior.
- Add comments for critical internal helpers when logic is non-obvious.

## Package and Naming

- Keep package names short and lowercase.
- Use clear type names for contracts and DTOs.
- Prefer interface-first design for transport/service/repository boundaries.

## Imports

- Use module-qualified local imports: gofoundry/<package>.
- Keep imports grouped by stdlib, internal module, external dependencies.

## Error Handling

- Return explicit errors; avoid silent failures.
- Use helper constructors for consistent API error semantics.

## Generics

- Use generics where they reduce duplication and preserve clarity.
- Document type parameters and expected constraints in comments.
