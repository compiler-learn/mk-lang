# Build
FROM prologic/go-builder:latest AS build

# Runtime
FROM alpine

COPY --from=build /src/monkey-lang /monkey-lang

ENTRYPOINT ["/monkey-lang"]
CMD [""]
