import { describe, it, expect } from "vitest";

describe("기본 테스트", () => {
  it("더하기 함수가 올바르게 작동한다", () => {
    const add = (a: number, b: number) => a + b;
    expect(add(2, 3)).toBe(5);
  });

  it("문자열 연결이 올바르게 작동한다", () => {
    const hello = "Hello";
    const world = "World";
    expect(`${hello} ${world}`).toBe("Hello World");
  });

  it("배열 길이 체크가 올바르게 작동한다", () => {
    const items = [1, 2, 3, 4, 5];
    expect(items).toHaveLength(5);
  });

  it("객체 속성 체크가 올바르게 작동한다", () => {
    const user = { name: "John", age: 30 };
    expect(user).toHaveProperty("name");
    expect(user.name).toBe("John");
  });

  it("비동기 함수가 올바르게 작동한다", async () => {
    const asyncFunction = async () => {
      await new Promise((resolve) => setTimeout(resolve, 10));
      return "async result";
    };

    const result = await asyncFunction();
    expect(result).toBe("async result");
  });
});
