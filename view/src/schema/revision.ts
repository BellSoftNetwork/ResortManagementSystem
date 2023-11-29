type RevisionType = {
  [key: string]: {
    name: string;
    color: string;
    icon: string;
  };
};

export const REVISION_TYPE_MAP: RevisionType = {
  CREATED: { name: "생성", color: "primary", icon: "add" } as const,
  UPDATED: { name: "변경", color: "warning", icon: "edit" } as const,
  DELETED: { name: "삭제", color: "red", icon: "remove" } as const,
} as const;
export type HistoryType = keyof typeof REVISION_TYPE_MAP;

export type Revision<T> = {
  entity: T;
  historyType: HistoryType;
  historyCreatedAt: string;
  updatedFields: string[];
};
