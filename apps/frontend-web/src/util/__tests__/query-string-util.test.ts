import { describe, it, expect } from "vitest";
import { formatSortParam } from "../query-string-util";
import type { SortRequest } from "src/schema/response";

describe("query-string-util", () => {
  describe("formatSortParam", () => {
    it("내림차순 정렬 파라미터를 포맷한다", () => {
      const sortRequest: SortRequest = {
        field: "createdAt",
        isDescending: true,
      };

      expect(formatSortParam(sortRequest)).toBe("createdAt,desc");
    });

    it("오름차순 정렬 파라미터를 포맷한다", () => {
      const sortRequest: SortRequest = {
        field: "name",
        isDescending: false,
      };

      expect(formatSortParam(sortRequest)).toBe("name,asc");
    });

    it("다양한 필드명을 처리한다", () => {
      expect(
        formatSortParam({ field: "id", isDescending: true })
      ).toBe("id,desc");
      
      expect(
        formatSortParam({ field: "updatedAt", isDescending: false })
      ).toBe("updatedAt,asc");
      
      expect(
        formatSortParam({ field: "price", isDescending: true })
      ).toBe("price,desc");
    });
  });
});
