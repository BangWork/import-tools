interface ImportMetaEnv {
  // more environment variable
  readonly VITE_PROXY_DOMAIN_REAL: string;
  readonly VITE_HOST: string
  readonly VITE_PORT: number
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
