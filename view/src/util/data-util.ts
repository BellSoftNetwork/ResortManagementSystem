import { EnumMap } from "src/types/map";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function isFormValueChanged(source: any, formData: any) {
  return Object.keys(formData).some((key) => source[key] !== formData[key]);
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function getPatchedFormData(source: any, formData: any) {
  const patchedData: EnumMap = {};

  if (source === undefined || formData === undefined) return patchedData;

  return Object.keys(formData)
    .filter((key: string) => {
      const sourceValue = source[key];
      const formValue = formData[key];
      const isObject = formValue !== null && typeof formValue === "object" && !Array.isArray(formValue);

      if (isObject && sourceValue.hasOwnProperty("id")) {
        return sourceValue.id !== formValue.id;
      } else {
        return sourceValue !== formValue;
      }
    })
    .reduce((result, key: string) => {
      result[key] = formData[key];
      return result;
    }, patchedData);
}
