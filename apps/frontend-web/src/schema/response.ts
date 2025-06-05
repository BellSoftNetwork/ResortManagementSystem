export type PageRequestParams = {
  page: number;
  size: number;
};

export type SortRequestParams = {
  sort: string;
};

export type SortRequest = {
  field: string;
  isDescending?: boolean;
};

type PageResponse = {
  index: number;
  size: number;
  totalElements: number;
  totalPages: number;
};

export type ListResponse<T> = {
  page: PageResponse;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  filter?: any;
  values: T[];
};

export type SingleResponse<T> = {
  value: T;
};
