export type ServerConfig = {
  isAvailableRegistration: boolean;
};

export type ServerEnv = {
  applicationFullName: string;
  applicationShortName: string;
  commitSha: string;
  commitShortSha: string;
  commitTitle: string;
  commitTimestamp: string;
};
