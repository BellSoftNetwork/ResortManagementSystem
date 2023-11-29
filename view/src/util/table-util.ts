import { FieldDetail } from "src/types/map"

type TableColumnDefinition = {
  label: string;
  name: string;
  field: string;
  format?: (value: string | number | null) => string;
};

export function convertTableColumnDef(
  fieldDetail: FieldDetail | null,
): TableColumnDefinition | null {
  return fieldDetail
    ? ({
      ...fieldDetail,
      name: fieldDetail.field,
    } as TableColumnDefinition)
    : null;
}
