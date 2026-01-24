import { ValidationRule } from "quasar";

export type EnumMap = {
  [field: string]: string;
};

export type FieldDetail = {
  label: string;
  field?: string;

  format?: (value: any) => any;
};

export type FieldMap = {
  [field: string]: FieldDetail;
};

export type StaticRuleMap = {
  [field: string]: ValidationRule[];
};

export type DynamicRuleMap = {
  [field: string]: (value: any) => ValidationRule[];
};
