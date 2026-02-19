import { ref, watch, Ref } from "vue";
import { useRoute, useRouter, LocationQuery } from "vue-router";
import { formatSortParam } from "src/util/query-string-util";
import { ListResponse } from "src/schema/response";

/**
 * Quasar QTable pagination object
 */
export interface QTablePagination {
  sortBy: string;
  descending: boolean;
  page: number;
  rowsPerPage: number;
  rowsNumber?: number;
}

/**
 * Quasar QTable onRequest props
 */
export interface QTableRequestProps {
  pagination: QTablePagination;
  filter?: any;
}

/**
 * API fetch function signature
 */
export type FetchFunction<T> = (params: {
  page: number;
  size: number;
  sort: string;
  [key: string]: any;
}) => Promise<ListResponse<T>>;

/**
 * useTable options
 */
export interface UseTableOptions<T> {
  /**
   * API fetch function
   */
  fetchFn: FetchFunction<T>;

  /**
   * Default pagination configuration
   */
  defaultPagination?: Partial<QTablePagination>;

  /**
   * Enable URL query string synchronization
   * @default true
   */
  syncUrl?: boolean;

  /**
   * Additional filter parameters
   */
  filter?: Ref<Record<string, any>>;

  /**
   * Error handler callback
   */
  onError?: (error: any) => void;
}

/**
 * useTable composable return type
 */
export interface UseTableReturn<T> {
  /**
   * Reactive pagination state (bind to QTable v-model:pagination)
   */
  pagination: Ref<QTablePagination>;

  /**
   * Reactive loading state (bind to QTable :loading)
   */
  loading: Ref<boolean>;

  /**
   * Reactive rows data (bind to QTable :rows)
   */
  rows: Ref<T[]>;

  /**
   * Total number of rows
   */
  totalRows: Ref<number>;

  /**
   * QTable onRequest handler
   */
  onRequest: (props: QTableRequestProps) => Promise<void>;

  /**
   * Reload current page
   */
  reload: () => void;

  /**
   * Reset to first page
   */
  resetToFirstPage: () => void;
}

/**
 * Reusable table composable for Quasar QTable with pagination, sorting, filtering, and URL sync
 *
 * @example
 * ```typescript
 * const { pagination, loading, rows, onRequest, reload } = useTable({
 *   fetchFn: fetchRooms,
 *   defaultPagination: {
 *     sortBy: 'number',
 *     descending: false,
 *     page: 1,
 *     rowsPerPage: 15,
 *   },
 * });
 * ```
 */
export function useTable<T>(options: UseTableOptions<T>): UseTableReturn<T> {
  const { fetchFn, defaultPagination = {}, syncUrl = true, filter, onError } = options;

  const route = useRoute();
  const router = useRouter();

  // Default configuration
  const defaultConfig: QTablePagination = {
    sortBy: defaultPagination.sortBy ?? "",
    descending: defaultPagination.descending ?? false,
    page: defaultPagination.page ?? 1,
    rowsPerPage: defaultPagination.rowsPerPage ?? 20,
    rowsNumber: 0,
  };

  // Reactive state
  const pagination = ref<QTablePagination>({ ...defaultConfig });
  const loading = ref(false);
  const rows = ref<T[]>([]);
  const totalRows = ref(0);

  // Internal ref for table component
  let tableRef: any = null;

  /**
   * Load pagination state from URL query string
   */
  function loadFromUrl() {
    if (!syncUrl) return;

    pagination.value.sortBy = route.query.sortBy?.toString() ?? defaultConfig.sortBy;
    pagination.value.descending = Boolean(route.query.descending ?? defaultConfig.descending);
    pagination.value.page = Number(route.query.page ?? defaultConfig.page);
    pagination.value.rowsPerPage = Number(route.query.rowsPerPage ?? defaultConfig.rowsPerPage);
  }

  /**
   * Sync pagination state to URL query string
   */
  function syncToUrl() {
    if (!syncUrl) return;

    const query: LocationQuery = {
      ...route.query,
    };

    // Only include non-default values
    if (pagination.value.page !== defaultConfig.page) {
      query.page = String(pagination.value.page);
    } else {
      delete query.page;
    }

    if (pagination.value.rowsPerPage !== defaultConfig.rowsPerPage) {
      query.rowsPerPage = String(pagination.value.rowsPerPage);
    } else {
      delete query.rowsPerPage;
    }

    if (pagination.value.sortBy !== defaultConfig.sortBy) {
      query.sortBy = pagination.value.sortBy;
    } else {
      delete query.sortBy;
    }

    if (pagination.value.descending !== defaultConfig.descending) {
      query.descending = String(pagination.value.descending);
    } else {
      delete query.descending;
    }

    router.push({ query });
  }

  /**
   * QTable onRequest handler
   */
  async function onRequest(props: QTableRequestProps): Promise<void> {
    const { page, rowsPerPage, sortBy, descending } = props.pagination;

    loading.value = true;

    try {
      const params: any = {
        page: page - 1, // API uses 0-based index
        size: rowsPerPage,
        sort: formatSortParam({ field: sortBy, isDescending: descending }),
      };

      // Merge additional filter parameters
      if (filter?.value) {
        Object.assign(params, filter.value);
      }

      const response = await fetchFn(params);

      // Update rows
      rows.value = response.values ?? [];
      totalRows.value = response.page.totalElements;

      // Update pagination
      pagination.value.rowsNumber = response.page.totalElements;
      pagination.value.page = response.page.index + 1; // Convert to 1-based
      pagination.value.rowsPerPage = response.page.size;
      pagination.value.sortBy = sortBy;
      pagination.value.descending = descending;

      // Sync to URL
      syncToUrl();
    } catch (error) {
      // Clear rows on error
      rows.value = [];
      totalRows.value = 0;

      // Call error handler if provided
      if (onError) {
        onError(error);
      } else {
        console.error("Table fetch error:", error);
      }
    } finally {
      loading.value = false;
    }
  }

  /**
   * Reload current page
   */
  function reload() {
    if (tableRef) {
      tableRef.requestServerInteraction();
    }
  }

  /**
   * Reset to first page
   */
  function resetToFirstPage() {
    pagination.value.page = 1;
    reload();
  }

  // Load initial state from URL
  loadFromUrl();

  // Watch route changes
  if (syncUrl) {
    watch(route, () => {
      loadFromUrl();
    });
  }

  return {
    pagination,
    loading,
    rows,
    totalRows,
    onRequest,
    reload,
    resetToFirstPage,
  };
}
