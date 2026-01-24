import { EnumMap } from "src/types/map";

export function isFormValueChanged(source: any, formData: any) {
  return Object.keys(formData).some((key) => source[key] !== formData[key]);
}

export function getPatchedFormData(source: any, formData: any) {
  const patchedData: EnumMap = {};

  if (source === undefined || formData === undefined) return patchedData;

  return Object.keys(formData)
    .filter((key: string) => {
      const sourceValue = source?.[key];
      const formValue = formData?.[key];
      const isObject = formValue !== null && typeof formValue === "object" && !Array.isArray(formValue);

      if (isObject && sourceValue?.["id"]) {
        return sourceValue.id !== formValue.id;
      } else if (Array.isArray(sourceValue) && Array.isArray(formValue)) {
        return !compareIdArrays(sourceValue, formValue);
      } else {
        return sourceValue !== formValue;
      }
    })
    .reduce((result, key: string) => {
      result[key] = formData[key];
      return result;
    }, patchedData);
}

function compareIdArrays(source: any[], target: any[]) {
  if (source.length !== target.length) return false;

  const sourceSet = new Set(source.map((item) => item.id));
  const targetSet = new Set(target.map((item) => item.id));

  return [...sourceSet].every((value) => targetSet.has(value));
}
