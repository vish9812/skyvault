import { get, handleJSONResponse } from "@sv/apis/common";
import { SystemConfig } from "./models";

export async function getSystemConfig(): Promise<SystemConfig> {
  const response = await get("system/config");
  return handleJSONResponse<SystemConfig>(response);
}

export * from "./models";