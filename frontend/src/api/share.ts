import api from "./index";
import type { ApiResponse, ShareResult } from "@/types";

export function getSharedPlaces({ username, group }: { username: string; group?: string }) {
  return api.get<ApiResponse<ShareResult>>(`/share/${username}`, {
    params: { group },
  });
}
