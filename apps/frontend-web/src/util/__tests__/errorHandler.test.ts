import { describe, it, expect, vi, beforeEach } from "vitest";

vi.unmock("axios");
import axios, { AxiosError, AxiosHeaders } from "axios";
import { getErrorMessage } from "../errorHandler";

describe("errorHandler", () => {
  describe("getErrorMessage", () => {
    it("Axios 에러가 아닌 경우 기본 메시지를 반환한다", () => {
      const error = new Error("Some error");

      expect(getErrorMessage(error)).toBe("요청에 실패했습니다");
    });

    it("null이나 undefined에 대해 기본 메시지를 반환한다", () => {
      expect(getErrorMessage(null)).toBe("요청에 실패했습니다");
      expect(getErrorMessage(undefined)).toBe("요청에 실패했습니다");
    });

    it("네트워크 에러 시 연결 실패 메시지를 반환한다", () => {
      const axiosError = new AxiosError(
        "Network Error",
        AxiosError.ERR_NETWORK,
        undefined,
        undefined,
        undefined
      );

      expect(getErrorMessage(axiosError)).toBe("서버에 연결할 수 없습니다");
    });

    it("500 서버 에러 시 서버 문제 메시지를 반환한다", () => {
      const axiosError = new AxiosError(
        "Server Error",
        AxiosError.ERR_BAD_RESPONSE,
        undefined,
        undefined,
        {
          status: 500,
          statusText: "Internal Server Error",
          data: {},
          headers: {},
          config: { headers: new AxiosHeaders() },
        }
      );

      expect(getErrorMessage(axiosError)).toBe(
        "서버에 문제가 발생했습니다 (오류 코드: 500)"
      );
    });

    it("502 서버 에러를 처리한다", () => {
      const axiosError = new AxiosError(
        "Bad Gateway",
        AxiosError.ERR_BAD_RESPONSE,
        undefined,
        undefined,
        {
          status: 502,
          statusText: "Bad Gateway",
          data: {},
          headers: {},
          config: { headers: new AxiosHeaders() },
        }
      );

      expect(getErrorMessage(axiosError)).toBe(
        "서버에 문제가 발생했습니다 (오류 코드: 502)"
      );
    });

    it("503 서버 에러를 처리한다", () => {
      const axiosError = new AxiosError(
        "Service Unavailable",
        AxiosError.ERR_BAD_RESPONSE,
        undefined,
        undefined,
        {
          status: 503,
          statusText: "Service Unavailable",
          data: {},
          headers: {},
          config: { headers: new AxiosHeaders() },
        }
      );

      expect(getErrorMessage(axiosError)).toBe(
        "서버에 문제가 발생했습니다 (오류 코드: 503)"
      );
    });

    it("400 에러 시 응답 메시지를 반환한다", () => {
      const axiosError = new AxiosError(
        "Bad Request",
        AxiosError.ERR_BAD_REQUEST,
        undefined,
        undefined,
        {
          status: 400,
          statusText: "Bad Request",
          data: { message: "잘못된 요청입니다" },
          headers: {},
          config: { headers: new AxiosHeaders() },
        }
      );

      expect(getErrorMessage(axiosError)).toBe("잘못된 요청입니다");
    });

    it("401 에러 시 응답 메시지를 반환한다", () => {
      const axiosError = new AxiosError(
        "Unauthorized",
        AxiosError.ERR_BAD_REQUEST,
        undefined,
        undefined,
        {
          status: 401,
          statusText: "Unauthorized",
          data: { message: "인증이 필요합니다" },
          headers: {},
          config: { headers: new AxiosHeaders() },
        }
      );

      expect(getErrorMessage(axiosError)).toBe("인증이 필요합니다");
    });

    it("응답 메시지가 없으면 기본 메시지를 반환한다", () => {
      const axiosError = new AxiosError(
        "Bad Request",
        AxiosError.ERR_BAD_REQUEST,
        undefined,
        undefined,
        {
          status: 400,
          statusText: "Bad Request",
          data: {},
          headers: {},
          config: { headers: new AxiosHeaders() },
        }
      );

      expect(getErrorMessage(axiosError)).toBe("요청에 실패했습니다");
    });

    it("빈 문자열 메시지는 기본 메시지를 반환한다", () => {
      const axiosError = new AxiosError(
        "Bad Request",
        AxiosError.ERR_BAD_REQUEST,
        undefined,
        undefined,
        {
          status: 400,
          statusText: "Bad Request",
          data: { message: "   " },
          headers: {},
          config: { headers: new AxiosHeaders() },
        }
      );

      expect(getErrorMessage(axiosError)).toBe("요청에 실패했습니다");
    });
  });
});
