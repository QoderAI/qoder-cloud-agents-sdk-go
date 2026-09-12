package convention

// packageVersion is reported to the server in User-Agent and
// X-Qoder-Package-Version. It must match the released module tag, otherwise the
// server-side SDK version statistics describe a version that was never shipped.
// `make check-version` compares it against the latest tag.
const packageVersion = "0.0.1-dev2"
