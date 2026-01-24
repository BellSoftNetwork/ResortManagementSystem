#!/usr/bin/env python3
"""
API 호환성 테스트 스크립트
api-legacy와 api-core의 응답을 비교하고 호환성을 검증합니다.

Usage:
    python3 scripts/api-compatibility-test.py [options]

Examples:
    python3 scripts/api-compatibility-test.py --core-only
    python3 scripts/api-compatibility-test.py --save-golden
    python3 scripts/api-compatibility-test.py --compare-golden
"""

import argparse
import json
import os
import sys
import time
import hashlib
import subprocess
import difflib
from pathlib import Path
from datetime import datetime
from typing import Dict, Any, List, Optional, Tuple
from dataclasses import dataclass, field, asdict
import requests
from requests.exceptions import ConnectionError, RequestException
import yaml


@dataclass
class EndpointConfig:
    id: str
    path: str
    method: str
    auth: bool = False
    role: Optional[str] = None
    description: str = ""
    params: Dict[str, Any] = field(default_factory=dict)
    path_params: Dict[str, Any] = field(default_factory=dict)
    request_body: Optional[Dict[str, Any]] = None
    status: str = "pending"
    expected_status: List[int] = field(default_factory=lambda: [200, 201, 204])


@dataclass
class TestResult:
    endpoint_id: str
    path: str
    method: str
    core_status: int = 0
    legacy_status: int = 0
    core_response: Optional[Dict] = None
    legacy_response: Optional[Dict] = None
    match: bool = False
    diff: List[str] = field(default_factory=list)
    error: Optional[str] = None
    duration_ms: float = 0


class APICompatibilityTester:
    def __init__(self, core_url: str = "http://localhost:8080", 
                 legacy_url: str = "http://localhost:8081",
                 manifest_path: str = None,
                 golden_dir: str = None):
        self.core_url = core_url.rstrip("/")
        self.legacy_url = legacy_url.rstrip("/")
        self.manifest_path = manifest_path or "tests/contracts/api-endpoints-manifest.yaml"
        self.golden_dir = Path(golden_dir or "tests/contracts/golden")
        self.golden_dir.mkdir(parents=True, exist_ok=True)
        
        self.tokens = {}
        self.test_accounts = {
            "SUPER_ADMIN": {"username": "testadmin", "password": "testadmin123"},
            "ADMIN": {"username": "testmanager", "password": "testmanager123"},
            "USER": {"username": "testuser", "password": "testuser123"}
        }
        
        self.ignore_fields = [
            "createdAt", "updatedAt", "deletedAt", "lastLoginAt", 
            "timestamp", "accessToken", "refreshToken", "accessTokenExpiresIn",
            "uptime", "hostname", "version"
        ]
        
        self.endpoints: List[EndpointConfig] = []
        self.results: List[TestResult] = []
        
    def load_manifest(self) -> bool:
        try:
            with open(self.manifest_path, 'r', encoding='utf-8') as f:
                manifest = yaml.safe_load(f)
            
            if manifest.get("ignore_fields"):
                self.ignore_fields = manifest["ignore_fields"]
            
            if manifest.get("test_accounts"):
                for role, account in manifest["test_accounts"].items():
                    self.test_accounts[role.upper()] = {
                        "username": account["username"],
                        "password": account["password"]
                    }
            
            for ep in manifest.get("endpoints", []):
                self.endpoints.append(EndpointConfig(
                    id=ep["id"],
                    path=ep["path"],
                    method=ep["method"],
                    auth=ep.get("auth", False),
                    role=ep.get("role"),
                    description=ep.get("description", ""),
                    params=ep.get("params", {}),
                    path_params=ep.get("path_params", {}),
                    request_body=ep.get("request_body"),
                    status=ep.get("status", "pending"),
                    expected_status=ep.get("expected_status", [200, 201, 204])
                ))
            
            return True
        except Exception as e:
            print(f"Failed to load manifest: {e}")
            return False
    
    def get_token(self, role: str, server: str = "core") -> Optional[str]:
        cache_key = f"{role}_{server}"
        if cache_key in self.tokens:
            return self.tokens[cache_key]
        
        account = self.test_accounts.get(role)
        if not account:
            return None
        
        base_url = self.core_url if server == "core" else self.legacy_url
        
        try:
            response = requests.post(
                f"{base_url}/api/v1/auth/login",
                json={"username": account["username"], "password": account["password"]},
                headers={"Content-Type": "application/json"},
                timeout=10
            )
            
            if response.status_code == 200:
                data = response.json()
                inner = data.get("value") or data.get("data") or data
                token = inner.get("accessToken")
                if token:
                    self.tokens[cache_key] = token
                    return token
        except Exception as e:
            print(f"Login failed for {role} on {server}: {e}")
        
        return None
    
    def setup_test_data(self, server: str = "core") -> bool:
        token = self.get_token("SUPER_ADMIN", server)
        if not token:
            print(f"Cannot get SUPER_ADMIN token for {server}")
            return False
        
        base_url = self.core_url if server == "core" else self.legacy_url
        
        try:
            response = requests.post(
                f"{base_url}/api/v1/dev/test-data",
                json={"type": "all"},
                headers={
                    "Authorization": f"Bearer {token}",
                    "Content-Type": "application/json"
                },
                timeout=30
            )
            return response.status_code in [200, 201]
        except Exception as e:
            print(f"Failed to setup test data: {e}")
            return False
    
    def make_request(self, endpoint: EndpointConfig, server: str = "core") -> Tuple[int, Optional[Dict]]:
        base_url = self.core_url if server == "core" else self.legacy_url
        
        path = endpoint.path
        for key, value in endpoint.path_params.items():
            path = path.replace(f"{{{key}}}", str(value))
        
        headers = {"Content-Type": "application/json"}
        
        if endpoint.auth:
            role = endpoint.role or "USER"
            token = self.get_token(role, server)
            if not token:
                return 401, {"error": f"Failed to get token for role {role}"}
            headers["Authorization"] = f"Bearer {token}"
        
        url = f"{base_url}{path}"
        
        try:
            response = requests.request(
                method=endpoint.method,
                url=url,
                headers=headers,
                params=endpoint.params if endpoint.method == "GET" else None,
                json=endpoint.request_body if endpoint.method in ["POST", "PATCH", "PUT"] else None,
                timeout=15
            )
            
            try:
                return response.status_code, response.json()
            except:
                return response.status_code, {"text": response.text}
        except ConnectionError:
            return 0, {"error": f"Connection refused: {url}"}
        except Exception as e:
            return 0, {"error": str(e)}
    
    def normalize_response(self, response: Optional[Dict]) -> Dict:
        if response is None:
            return {}
        
        def _clean(obj):
            if isinstance(obj, dict):
                return {k: _clean(v) for k, v in obj.items() if k not in self.ignore_fields}
            elif isinstance(obj, list):
                return [_clean(item) for item in obj]
            return obj
        
        return _clean(response)
    
    def compare_responses(self, legacy: Optional[Dict], core: Optional[Dict]) -> Tuple[bool, List[str]]:
        legacy_norm = self.normalize_response(legacy)
        core_norm = self.normalize_response(core)
        
        if legacy_norm == core_norm:
            return True, []
        
        legacy_json = json.dumps(legacy_norm, indent=2, sort_keys=True, ensure_ascii=False)
        core_json = json.dumps(core_norm, indent=2, sort_keys=True, ensure_ascii=False)
        
        diff = list(difflib.unified_diff(
            legacy_json.splitlines(keepends=True),
            core_json.splitlines(keepends=True),
            fromfile="api-legacy",
            tofile="api-core",
            lineterm=""
        ))
        
        return False, diff
    
    def save_golden(self, endpoint: EndpointConfig, response: Dict):
        filename = f"{endpoint.id}.json"
        filepath = self.golden_dir / filename
        
        normalized = self.normalize_response(response)
        with open(filepath, 'w', encoding='utf-8') as f:
            json.dump(normalized, f, indent=2, ensure_ascii=False)
    
    def load_golden(self, endpoint_id: str) -> Optional[Dict]:
        filepath = self.golden_dir / f"{endpoint_id}.json"
        if not filepath.exists():
            return None
        
        with open(filepath, 'r', encoding='utf-8') as f:
            return json.load(f)
    
    def test_endpoint(self, endpoint: EndpointConfig, mode: str = "compare") -> TestResult:
        start_time = time.time()
        result = TestResult(
            endpoint_id=endpoint.id,
            path=endpoint.path,
            method=endpoint.method
        )
        
        if endpoint.status in ["core-only", "legacy-only"]:
            server = "core" if endpoint.status == "core-only" else "legacy"
            status, response = self.make_request(endpoint, server)
            
            if server == "core":
                result.core_status = status
                result.core_response = response
            else:
                result.legacy_status = status
                result.legacy_response = response
            
            result.match = status in endpoint.expected_status
            result.duration_ms = (time.time() - start_time) * 1000
            return result
        
        core_status, core_response = self.make_request(endpoint, "core")
        result.core_status = core_status
        result.core_response = core_response
        
        if mode == "core-only":
            result.match = core_status in endpoint.expected_status
            result.duration_ms = (time.time() - start_time) * 1000
            return result
        
        if mode == "save-golden":
            if core_status in [200, 201]:
                self.save_golden(endpoint, core_response)
                result.match = True
            else:
                result.error = f"Cannot save golden: status {core_status}"
            result.duration_ms = (time.time() - start_time) * 1000
            return result
        
        if mode == "compare-golden":
            golden = self.load_golden(endpoint.id)
            if golden is None:
                result.error = "No golden file found"
            else:
                result.match, result.diff = self.compare_responses(golden, core_response)
            result.duration_ms = (time.time() - start_time) * 1000
            return result
        
        legacy_status, legacy_response = self.make_request(endpoint, "legacy")
        result.legacy_status = legacy_status
        result.legacy_response = legacy_response
        
        if core_status != legacy_status:
            result.error = f"Status mismatch: core={core_status}, legacy={legacy_status}"
        else:
            result.match, result.diff = self.compare_responses(legacy_response, core_response)
        
        result.duration_ms = (time.time() - start_time) * 1000
        return result
    
    def run_tests(self, mode: str = "compare", filter_ids: List[str] = None) -> List[TestResult]:
        self.results = []
        
        endpoints = self.endpoints
        if filter_ids:
            endpoints = [ep for ep in self.endpoints if ep.id in filter_ids]
        
        for i, endpoint in enumerate(endpoints):
            print(f"[{i+1}/{len(endpoints)}] Testing {endpoint.method} {endpoint.path}...", end=" ")
            
            result = self.test_endpoint(endpoint, mode)
            self.results.append(result)
            
            if result.error:
                print(f"ERROR: {result.error}")
            elif result.match:
                print(f"OK ({result.duration_ms:.0f}ms)")
            else:
                print(f"MISMATCH ({result.duration_ms:.0f}ms)")
        
        return self.results
    
    def print_summary(self):
        total = len(self.results)
        passed = sum(1 for r in self.results if r.match)
        failed = sum(1 for r in self.results if not r.match and not r.error)
        errors = sum(1 for r in self.results if r.error)
        
        print("\n" + "=" * 60)
        print("API Compatibility Test Summary")
        print("=" * 60)
        print(f"Total: {total}")
        print(f"  Passed:  {passed}")
        print(f"  Failed:  {failed}")
        print(f"  Errors:  {errors}")
        print("=" * 60)
        
        if failed > 0:
            print("\nFailed Endpoints:")
            for r in self.results:
                if not r.match and not r.error:
                    print(f"  - {r.method} {r.path} ({r.endpoint_id})")
                    if r.diff and len(r.diff) <= 10:
                        for line in r.diff[:10]:
                            print(f"      {line.rstrip()}")
        
        if errors > 0:
            print("\nError Endpoints:")
            for r in self.results:
                if r.error:
                    print(f"  - {r.method} {r.path}: {r.error}")
    
    def save_results(self, output_path: str = "api-compatibility-results.json"):
        results_data = {
            "timestamp": datetime.now().isoformat(),
            "core_url": self.core_url,
            "legacy_url": self.legacy_url,
            "total": len(self.results),
            "passed": sum(1 for r in self.results if r.match),
            "failed": sum(1 for r in self.results if not r.match and not r.error),
            "errors": sum(1 for r in self.results if r.error),
            "results": [
                {
                    "id": r.endpoint_id,
                    "path": r.path,
                    "method": r.method,
                    "match": r.match,
                    "core_status": r.core_status,
                    "legacy_status": r.legacy_status,
                    "error": r.error,
                    "diff_lines": len(r.diff) if r.diff else 0,
                    "duration_ms": r.duration_ms
                }
                for r in self.results
            ]
        }
        
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(results_data, f, indent=2, ensure_ascii=False)
        
        print(f"\nResults saved to {output_path}")


def main():
    parser = argparse.ArgumentParser(
        description="API Compatibility Tester for Resort Management System",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument("--core-url", default="http://localhost:8080")
    parser.add_argument("--legacy-url", default="http://localhost:8081")
    parser.add_argument("--manifest", default="tests/contracts/api-endpoints-manifest.yaml")
    parser.add_argument("--golden-dir", default="tests/contracts/golden")
    parser.add_argument("--output", default="api-compatibility-results.json")
    
    mode_group = parser.add_mutually_exclusive_group()
    mode_group.add_argument("--core-only", action="store_true", help="Test api-core only")
    mode_group.add_argument("--save-golden", action="store_true", help="Save api-core responses as golden files")
    mode_group.add_argument("--compare-golden", action="store_true", help="Compare api-core with golden files")
    
    parser.add_argument("--filter", nargs="*", help="Test specific endpoint IDs only")
    parser.add_argument("--setup-data", action="store_true", help="Setup test data before testing")
    parser.add_argument("-v", "--verbose", action="store_true")
    
    args = parser.parse_args()
    
    tester = APICompatibilityTester(
        core_url=args.core_url,
        legacy_url=args.legacy_url,
        manifest_path=args.manifest,
        golden_dir=args.golden_dir
    )
    
    if not tester.load_manifest():
        print("Failed to load manifest. Creating default endpoints...")
        tester.endpoints = [
            EndpointConfig("health", "/actuator/health", "GET", description="Health check"),
            EndpointConfig("env", "/api/v1/env", "GET", description="Environment info"),
            EndpointConfig("config", "/api/v1/config", "GET", description="Config info"),
        ]
    
    print(f"Loaded {len(tester.endpoints)} endpoints from manifest")
    
    if args.setup_data:
        print("Setting up test data...")
        if tester.setup_test_data("core"):
            print("Test data created successfully")
        else:
            print("Warning: Failed to create test data")
    
    mode = "compare"
    if args.core_only:
        mode = "core-only"
    elif args.save_golden:
        mode = "save-golden"
    elif args.compare_golden:
        mode = "compare-golden"
    
    print(f"\nRunning tests in '{mode}' mode...")
    print("-" * 60)
    
    tester.run_tests(mode=mode, filter_ids=args.filter)
    tester.print_summary()
    tester.save_results(args.output)
    
    failed_count = sum(1 for r in tester.results if not r.match)
    sys.exit(1 if failed_count > 0 else 0)


if __name__ == "__main__":
    main()
