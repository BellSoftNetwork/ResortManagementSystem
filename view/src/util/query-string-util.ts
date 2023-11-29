import { SortRequest } from "src/schema/response"

export function formatSortParam(sortRequest: SortRequest) {
  const order = sortRequest.isDescending ? "desc" : "asc"
  return `${sortRequest.field},${order}`
}
