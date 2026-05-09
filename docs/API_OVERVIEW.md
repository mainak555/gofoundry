# API Overview

This page summarizes key extension points and public contracts.

## Authentication

- auth.GetJWT
- auth.GetClaims
- auth.ValidateJWTIssuer
- auth.ValidateJWTSignature
- http/oidc.ValidateToken middleware

## Caching

- cache.ICache, cache.IMongoCache, cache.IRedisCache
- cache.NewRedisCache, cache.NewMongoCache, cache.NewCache
- cache.Get, cache.GetMany, cache.Set, cache.SetJson, cache.Delete

## MongoDB and Repository

- db/mongodb.IMongoClient and MongoClient lifecycle methods.
- repository/interfaces.IMongoRepository and TMongoRepository[T].
- repository.NewMongoRepository and repository.NewMongoTRepository[T].
- generics.ICommonService[T] and generics.NewCommonService[T].

## HTTP Transport

- http/server decode and encode helpers.
- http/chi decode helpers and middleware package.
- generics.ApiRouter fluent route composition methods.

## DTO and Helper Utilities

- dtos pagination/filter/sort request structures.
- helpers query and context accessor helpers.
- util shared conversion and reflection utilities.

Refer to package source for full signatures until example-driven docs are added.
