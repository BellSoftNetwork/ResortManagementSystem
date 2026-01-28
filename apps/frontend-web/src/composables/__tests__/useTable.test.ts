import { describe, it, expect, vi, beforeEach } from "vitest";
import { ref } from "vue";
import { useTable } from "../useTable";
import type { FetchFunction, QTableRequestProps } from "../useTable";
import type { ListResponse } from "src/schema/response";

// Mock vue-router
const mockPush = vi.fn();
const mockRoute = {
  query: {},
};

vi.mock("vue-router", () => ({
  useRoute: () => mockRoute,
  useRouter: () => ({
    push: mockPush,
  }),
}));

// Mock formatSortParam
vi.mock("src/util/query-string-util", () => ({
  formatSortParam: ({ field, isDescending }: { field: string; isDescending: boolean }) => {
    if (!field) return "";
    return isDescending ? `${field},desc` : `${field},asc`;
  },
}));

describe("useTable", () => {
  let mockFetchFn: ReturnType<typeof vi.fn<Parameters<FetchFunction<any>>, ReturnType<FetchFunction<any>>>>;

  beforeEach(() => {
    vi.clearAllMocks();
    mockRoute.query = {};
    mockFetchFn = vi.fn();
  });

  describe("초기화", () => {
    it("기본 pagination 값으로 초기화된다", () => {
      const { pagination } = useTable({ fetchFn: mockFetchFn });

      expect(pagination.value.page).toBe(1);
      expect(pagination.value.rowsPerPage).toBe(20);
      expect(pagination.value.sortBy).toBe("");
      expect(pagination.value.descending).toBe(false);
      expect(pagination.value.rowsNumber).toBe(0);
    });

    it("커스텀 pagination 값으로 초기화된다", () => {
      const { pagination } = useTable({
        fetchFn: mockFetchFn,
        defaultPagination: {
          page: 2,
          rowsPerPage: 15,
          sortBy: "name",
          descending: true,
        },
      });

      expect(pagination.value.page).toBe(2);
      expect(pagination.value.rowsPerPage).toBe(15);
      expect(pagination.value.sortBy).toBe("name");
      expect(pagination.value.descending).toBe(true);
    });

    it("초기 loading 상태는 false이다", () => {
      const { loading } = useTable({ fetchFn: mockFetchFn });

      expect(loading.value).toBe(false);
    });

    it("초기 rows는 빈 배열이다", () => {
      const { rows } = useTable({ fetchFn: mockFetchFn });

      expect(rows.value).toEqual([]);
    });

    it("초기 totalRows는 0이다", () => {
      const { totalRows } = useTable({ fetchFn: mockFetchFn });

      expect(totalRows.value).toBe(0);
    });
  });

  describe("onRequest 핸들러", () => {
    it("성공 시 rows와 pagination을 업데이트한다", async () => {
      const mockResponse: ListResponse<{ id: number; name: string }> = {
        values: [
          { id: 1, name: "Item 1" },
          { id: 2, name: "Item 2" },
        ],
        page: {
          totalElements: 100,
          index: 0, // 0-based
          size: 20,
        },
      };

      mockFetchFn.mockResolvedValue(mockResponse);

      const { onRequest, rows, totalRows, pagination } = useTable({
        fetchFn: mockFetchFn,
      });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "name",
          descending: false,
        },
      };

      await onRequest(requestProps);

      expect(rows.value).toEqual(mockResponse.values);
      expect(totalRows.value).toBe(100);
      expect(pagination.value.rowsNumber).toBe(100);
      expect(pagination.value.page).toBe(1); // API는 0-based, UI는 1-based
      expect(pagination.value.rowsPerPage).toBe(20);
      expect(pagination.value.sortBy).toBe("name");
      expect(pagination.value.descending).toBe(false);
    });

    it("fetchFn을 올바른 파라미터로 호출한다", async () => {
      const mockResponse: ListResponse<any> = {
        values: [],
        page: { totalElements: 0, index: 0, size: 20 },
      };

      mockFetchFn.mockResolvedValue(mockResponse);

      const { onRequest } = useTable({ fetchFn: mockFetchFn });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 2,
          rowsPerPage: 15,
          sortBy: "createdAt",
          descending: true,
        },
      };

      await onRequest(requestProps);

      expect(mockFetchFn).toHaveBeenCalledWith({
        page: 1, // 2 - 1 (0-based)
        size: 15,
        sort: "createdAt,desc",
      });
    });

    it("정렬 없이 호출 시 빈 sort 파라미터를 전달한다", async () => {
      const mockResponse: ListResponse<any> = {
        values: [],
        page: { totalElements: 0, index: 0, size: 20 },
      };

      mockFetchFn.mockResolvedValue(mockResponse);

      const { onRequest } = useTable({ fetchFn: mockFetchFn });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "",
          descending: false,
        },
      };

      await onRequest(requestProps);

      expect(mockFetchFn).toHaveBeenCalledWith({
        page: 0,
        size: 20,
        sort: "",
      });
    });

    it("추가 필터 파라미터를 병합한다", async () => {
      const mockResponse: ListResponse<any> = {
        values: [],
        page: { totalElements: 0, index: 0, size: 20 },
      };

      mockFetchFn.mockResolvedValue(mockResponse);

      const filter = ref({
        status: "ACTIVE",
        type: "STANDARD",
      });

      const { onRequest } = useTable({
        fetchFn: mockFetchFn,
        filter,
      });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "",
          descending: false,
        },
      };

      await onRequest(requestProps);

      expect(mockFetchFn).toHaveBeenCalledWith({
        page: 0,
        size: 20,
        sort: "",
        status: "ACTIVE",
        type: "STANDARD",
      });
    });

    it("loading 상태를 올바르게 관리한다", async () => {
      const mockResponse: ListResponse<any> = {
        values: [],
        page: { totalElements: 0, index: 0, size: 20 },
      };

      let resolvePromise: (value: ListResponse<any>) => void;
      const promise = new Promise<ListResponse<any>>((resolve) => {
        resolvePromise = resolve;
      });

      mockFetchFn.mockReturnValue(promise);

      const { onRequest, loading } = useTable({ fetchFn: mockFetchFn });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "",
          descending: false,
        },
      };

      expect(loading.value).toBe(false);

      const requestPromise = onRequest(requestProps);

      // 요청 중에는 loading이 true
      expect(loading.value).toBe(true);

      resolvePromise!(mockResponse);
      await requestPromise;

      // 요청 완료 후에는 loading이 false
      expect(loading.value).toBe(false);
    });
  });

  describe("에러 처리", () => {
    it("에러 발생 시 rows를 초기화한다", async () => {
      const error = new Error("Network error");
      mockFetchFn.mockRejectedValue(error);

      const { onRequest, rows, totalRows } = useTable({
        fetchFn: mockFetchFn,
      });

      // 먼저 데이터를 채움
      rows.value = [{ id: 1 }];
      totalRows.value = 100;

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "",
          descending: false,
        },
      };

      await onRequest(requestProps);

      expect(rows.value).toEqual([]);
      expect(totalRows.value).toBe(0);
    });

    it("에러 발생 시 onError 콜백을 호출한다", async () => {
      const error = new Error("Network error");
      mockFetchFn.mockRejectedValue(error);

      const onError = vi.fn();

      const { onRequest } = useTable({
        fetchFn: mockFetchFn,
        onError,
      });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "",
          descending: false,
        },
      };

      await onRequest(requestProps);

      expect(onError).toHaveBeenCalledWith(error);
    });

    it("onError 콜백이 없으면 console.error를 호출한다", async () => {
      const error = new Error("Network error");
      mockFetchFn.mockRejectedValue(error);

      const consoleErrorSpy = vi.spyOn(console, "error").mockImplementation(() => {});

      const { onRequest } = useTable({
        fetchFn: mockFetchFn,
      });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "",
          descending: false,
        },
      };

      await onRequest(requestProps);

      expect(consoleErrorSpy).toHaveBeenCalledWith("Table fetch error:", error);

      consoleErrorSpy.mockRestore();
    });

    it("에러 발생 후에도 loading 상태를 false로 변경한다", async () => {
      const error = new Error("Network error");
      mockFetchFn.mockRejectedValue(error);

      const { onRequest, loading } = useTable({
        fetchFn: mockFetchFn,
      });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "",
          descending: false,
        },
      };

      await onRequest(requestProps);

      expect(loading.value).toBe(false);
    });
  });

  describe("resetToFirstPage", () => {
    it("페이지를 1로 리셋한다", () => {
      const { pagination, resetToFirstPage } = useTable({
        fetchFn: mockFetchFn,
        defaultPagination: { page: 5 },
      });

      expect(pagination.value.page).toBe(5);

      resetToFirstPage();

      expect(pagination.value.page).toBe(1);
    });
  });

  describe("URL 동기화", () => {
    it("syncUrl이 false일 때 URL을 동기화하지 않는다", async () => {
      const mockResponse: ListResponse<any> = {
        values: [],
        page: { totalElements: 0, index: 0, size: 20 },
      };

      mockFetchFn.mockResolvedValue(mockResponse);

      const { onRequest } = useTable({
        fetchFn: mockFetchFn,
        syncUrl: false,
      });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 2,
          rowsPerPage: 15,
          sortBy: "name",
          descending: true,
        },
      };

      await onRequest(requestProps);

      // router.push가 호출되지 않아야 함
      expect(mockPush).not.toHaveBeenCalled();
    });

    it("syncUrl이 true일 때 URL을 동기화한다", async () => {
      const mockResponse: ListResponse<any> = {
        values: [],
        page: { totalElements: 0, index: 1, size: 15 }, // API returns 0-based index
      };

      mockFetchFn.mockResolvedValue(mockResponse);

      const { onRequest } = useTable({
        fetchFn: mockFetchFn,
        syncUrl: true,
        defaultPagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "",
          descending: false,
        },
      });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 2,
          rowsPerPage: 15,
          sortBy: "name",
          descending: true,
        },
      };

      await onRequest(requestProps);

      expect(mockPush).toHaveBeenCalledWith({
        query: {
          page: "2",
          rowsPerPage: "15",
          sortBy: "name",
          descending: "true",
        },
      });
    });

    it("기본값과 동일한 파라미터는 URL에서 제거한다", async () => {
      const mockResponse: ListResponse<any> = {
        values: [],
        page: { totalElements: 0, index: 0, size: 20 },
      };

      mockFetchFn.mockResolvedValue(mockResponse);

      const { onRequest } = useTable({
        fetchFn: mockFetchFn,
        syncUrl: true,
        defaultPagination: {
          page: 1,
          rowsPerPage: 20,
          sortBy: "",
          descending: false,
        },
      });

      const requestProps: QTableRequestProps = {
        pagination: {
          page: 1, // 기본값
          rowsPerPage: 20, // 기본값
          sortBy: "", // 기본값
          descending: false, // 기본값
        },
      };

      await onRequest(requestProps);

      // 기본값이므로 query가 비어있어야 함
      expect(mockPush).toHaveBeenCalledWith({
        query: {},
      });
    });
  });

  describe("URL에서 초기 상태 로드", () => {
    it("URL query에서 pagination 상태를 로드한다", () => {
      mockRoute.query = {
        page: "3",
        rowsPerPage: "15",
        sortBy: "createdAt",
        descending: "true",
      };

      const { pagination } = useTable({
        fetchFn: mockFetchFn,
        syncUrl: true,
      });

      expect(pagination.value.page).toBe(3);
      expect(pagination.value.rowsPerPage).toBe(15);
      expect(pagination.value.sortBy).toBe("createdAt");
      expect(pagination.value.descending).toBe(true);
    });

    it("syncUrl이 false일 때 URL에서 로드하지 않는다", () => {
      mockRoute.query = {
        page: "3",
        rowsPerPage: "15",
      };

      const { pagination } = useTable({
        fetchFn: mockFetchFn,
        syncUrl: false,
        defaultPagination: {
          page: 1,
          rowsPerPage: 20,
        },
      });

      // 기본값 유지
      expect(pagination.value.page).toBe(1);
      expect(pagination.value.rowsPerPage).toBe(20);
    });
  });
});
