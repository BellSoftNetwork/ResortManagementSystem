#!/usr/bin/env python3
"""
Resort Management System API Test Script

This script provides authentication and API testing functionality for local development.
It automatically handles token management, user creation, and API calls.

Usage:
    python api-test.py <endpoint> [options]

Examples:
    python api-test.py /api/v1/users --method GET
    python api-test.py /api/v1/dev/generate-essential-data --method POST --role SUPER_ADMIN
    python api-test.py /api/v1/reservations --method POST --data '{"roomId": 1, "checkIn": "2024-01-01"}'
"""

import argparse
import json
import os
import sys
import time
import subprocess
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, Optional, Tuple
import requests
from requests.exceptions import ConnectionError, RequestException


class APITestClient:
    """API Test Client for Resort Management System"""
    
    def __init__(self, base_url: str = "http://localhost:8080", token_dir: str = ".tokens"):
        self.base_url = base_url.rstrip("/")
        self.token_dir = Path(token_dir)
        self.token_dir.mkdir(exist_ok=True)
        
        # Default test accounts by role
        self.test_accounts = {
            "SUPER_ADMIN": {
                "username": "testadmin",
                "password": "testadmin123"
            },
            "ADMIN": {
                "username": "testmanager",
                "password": "testmanager123"
            },
            "USER": {
                "username": "testuser",
                "password": "testuser123"
            }
        }
        
        # Pre-encoded BCrypt passwords for emergency recovery
        # These match the passwords above with {bcrypt} prefix for Spring Security compatibility
        self.encoded_passwords = {
            "testadmin123": "{bcrypt}$2a$10$dXJ3SW6G7P50lGmMkkmwe.20cQQubK3.HZWzG3YB1tlRy.fqvM/BG",
            "testmanager123": "{bcrypt}$2a$10$dXJ3SW6G7P50lGmMkkmwe.20cQQubK3.HZWzG3YB1tlRy.fqvM/BG",
            "testuser123": "{bcrypt}$2a$10$dXJ3SW6G7P50lGmMkkmwe.20cQQubK3.HZWzG3YB1tlRy.fqvM/BG"
        }
        
        # MySQL connection details for emergency recovery
        self.mysql_config = {
            "host": os.getenv("DATABASE_MYSQL_HOST", "localhost"),
            "port": os.getenv("DATABASE_MYSQL_PORT", "3306"),
            "user": os.getenv("DATABASE_MYSQL_USER", "root"),
            "password": os.getenv("DATABASE_MYSQL_PASSWORD", "root"),
            "database": os.getenv("DATABASE_MYSQL_DATABASE", "resort_management")
        }
    
    def _get_token_file(self, role: str) -> Path:
        """Get token file path for a specific role"""
        return self.token_dir / f"token_{role.lower()}.json"
    
    def _save_tokens(self, role: str, access_token: str, refresh_token: str, expires_in: int):
        """Save tokens to file"""
        token_data = {
            "access_token": access_token,
            "refresh_token": refresh_token,
            "expires_at": (datetime.now() + timedelta(seconds=expires_in - 60)).isoformat(),
            "created_at": datetime.now().isoformat()
        }
        
        token_file = self._get_token_file(role)
        with open(token_file, 'w') as f:
            json.dump(token_data, f, indent=2)
    
    def _load_tokens(self, role: str) -> Optional[Dict]:
        """Load tokens from file"""
        token_file = self._get_token_file(role)
        if not token_file.exists():
            return None
        
        try:
            with open(token_file, 'r') as f:
                return json.load(f)
        except:
            return None
    
    def _is_token_valid(self, token_data: Dict) -> bool:
        """Check if access token is still valid"""
        if not token_data or "expires_at" not in token_data:
            return False
        
        expires_at = datetime.fromisoformat(token_data["expires_at"])
        return datetime.now() < expires_at
    
    def _clear_login_attempts(self, username: str) -> bool:
        """Clear login attempts from database to unlock account"""
        try:
            # Use docker compose to execute MySQL command
            cmd = [
                "docker", "compose", "exec", "-T", "mysql",
                "mysql", "-u", self.mysql_config["user"], f"-p{self.mysql_config['password']}",
                self.mysql_config["database"], "-e",
                f"DELETE FROM login_attempts WHERE username = '{username}';"
            ]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            if result.returncode == 0:
                print(f"Cleared login attempts for {username}")
                return True
            else:
                # Try direct MySQL connection as fallback
                cmd = [
                    "mysql", "-h", self.mysql_config["host"], "-P", self.mysql_config["port"],
                    "-u", self.mysql_config["user"], f"-p{self.mysql_config['password']}",
                    self.mysql_config["database"], "-e",
                    f"DELETE FROM login_attempts WHERE username = '{username}';"
                ]
                result = subprocess.run(cmd, capture_output=True, text=True)
                return result.returncode == 0
        except Exception as e:
            print(f"Failed to clear login attempts: {e}")
            return False
    
    def _reset_user_password(self, username: str, role: str) -> bool:
        """Reset user password in database"""
        try:
            account = self.test_accounts.get(role)
            if not account:
                return False
            
            encoded_password = self.encoded_passwords.get(account["password"])
            if not encoded_password:
                return False
            
            # Use docker compose to execute MySQL command
            cmd = [
                "docker", "compose", "exec", "-T", "mysql",
                "mysql", "-u", self.mysql_config["user"], f"-p{self.mysql_config['password']}",
                self.mysql_config["database"], "-e",
                f"UPDATE users SET password = '{encoded_password}' WHERE user_id = '{username}';"
            ]
            
            result = subprocess.run(cmd, capture_output=True, text=True)
            if result.returncode == 0:
                print(f"Reset password for {username}")
                return True
            else:
                # Try direct MySQL connection as fallback
                cmd = [
                    "mysql", "-h", self.mysql_config["host"], "-P", self.mysql_config["port"],
                    "-u", self.mysql_config["user"], f"-p{self.mysql_config['password']}",
                    self.mysql_config["database"], "-e",
                    f"UPDATE users SET password = '{encoded_password}' WHERE user_id = '{username}';"
                ]
                result = subprocess.run(cmd, capture_output=True, text=True)
                return result.returncode == 0
        except Exception as e:
            print(f"Failed to reset password: {e}")
            return False
    
    def _refresh_access_token(self, role: str, refresh_token: str) -> Optional[str]:
        """Refresh access token using refresh token"""
        try:
            response = requests.post(
                f"{self.base_url}/api/v1/auth/refresh",
                json={"refreshToken": refresh_token},
                headers={"Content-Type": "application/json"}
            )
            
            if response.status_code == 200:
                data = response.json()
                value = data.get("value", data)  # Handle both wrapped and unwrapped responses
                
                # Calculate expires_in from timestamp if needed
                expires_in = 900  # Default 15 minutes
                if "accessTokenExpiresIn" in value:
                    # Convert timestamp to seconds
                    expires_in = int((value["accessTokenExpiresIn"] - time.time() * 1000) / 1000)
                elif "expiresIn" in value:
                    expires_in = value["expiresIn"]
                
                self._save_tokens(
                    role,
                    value["accessToken"],
                    value["refreshToken"],
                    expires_in
                )
                return value["accessToken"]
        except:
            pass
        
        return None
    
    def _login(self, username: str, password: str, role: str, retry_count: int = 0) -> Optional[str]:
        """Login and get access token with retry logic"""
        max_retries = 3
        
        try:
            response = requests.post(
                f"{self.base_url}/api/v1/auth/login",
                json={"username": username, "password": password},
                headers={"Content-Type": "application/json"}
            )
            
            if response.status_code == 200:
                data = response.json()
                value = data.get("value", data)  # Handle both wrapped and unwrapped responses
                
                # Calculate expires_in from timestamp if needed
                expires_in = 900  # Default 15 minutes
                if "accessTokenExpiresIn" in value:
                    # Convert timestamp to seconds
                    expires_in = int((value["accessTokenExpiresIn"] - time.time() * 1000) / 1000)
                elif "expiresIn" in value:
                    expires_in = value["expiresIn"]
                
                self._save_tokens(
                    role,
                    value["accessToken"],
                    value["refreshToken"],
                    expires_in
                )
                return value["accessToken"]
            
            elif response.status_code == 429:
                # Too many requests - account is locked
                print(f"Account locked due to too many login attempts")
                if retry_count < max_retries:
                    print("Attempting to unlock account...")
                    if self._clear_login_attempts(username):
                        time.sleep(2)  # Wait a bit before retrying
                        return self._login(username, password, role, retry_count + 1)
            
            elif response.status_code == 401:
                print(f"Login failed: Invalid credentials")
                if retry_count < max_retries:
                    print("Attempting to reset password...")
                    if self._reset_user_password(username, role):
                        time.sleep(1)
                        return self._login(username, password, role, retry_count + 1)
            
            else:
                print(f"Login failed: {response.status_code} - {response.text}")
                
        except ConnectionError:
            print(f"Failed to connect to API at {self.base_url}")
            print("Make sure the API server is running with: docker compose up -d api-core")
        except Exception as e:
            print(f"Login error: {e}")
        
        return None
    
    def _register(self, username: str, password: str, role: str) -> bool:
        """Register a new user"""
        try:
            # For admin/super_admin registration, we need an existing admin token
            if role != "USER":
                # Try to get super admin token first
                admin_token = None
                for admin_role in ["SUPER_ADMIN", "ADMIN"]:
                    if admin_role != role:  # Don't use the same role we're trying to create
                        admin_token = self.get_access_token(admin_role)
                        if admin_token:
                            break
                
                if not admin_token:
                    print(f"Cannot create {role} user - no admin access available")
                    return False
                
                # Use admin endpoint to create user with specific role
                user_data = {
                    "userId": username,
                    "username": username,
                    "password": password,
                    "name": f"Test {role.replace('_', ' ').title()}",
                    "email": f"{username}@test.com",
                    "phoneNumber": "010-0000-0000",
                    "authorities": role
                }
                
                response = requests.post(
                    f"{self.base_url}/api/v1/admin/accounts",
                    json=user_data,
                    headers={
                        "Content-Type": "application/json",
                        "Authorization": f"Bearer {admin_token}"
                    }
                )
            else:
                # Regular user registration
                user_data = {
                    "userId": username,  # Add userId field
                    "username": username,
                    "password": password,
                    "name": f"Test User",
                    "email": f"{username}@test.com",
                    "phoneNumber": "010-0000-0000"
                }
                
                response = requests.post(
                    f"{self.base_url}/api/v1/auth/register",
                    json=user_data,
                    headers={"Content-Type": "application/json"}
                )
            
            if response.status_code in [200, 201]:
                print(f"Successfully created user: {username} with role: {role}")
                return True
            elif response.status_code == 409:
                print(f"User {username} already exists")
                # Try to reset the password instead
                return self._reset_user_password(username, role)
            else:
                print(f"Registration failed: {response.status_code} - {response.text}")
                
        except Exception as e:
            print(f"Registration error: {e}")
        
        return False
    
    def get_access_token(self, role: str = "USER") -> Optional[str]:
        """Get valid access token for the specified role"""
        # Try to load existing token
        token_data = self._load_tokens(role)
        
        # If token is valid, return it
        if token_data and self._is_token_valid(token_data):
            return token_data["access_token"]
        
        # If we have refresh token, try to refresh
        if token_data and "refresh_token" in token_data:
            access_token = self._refresh_access_token(role, token_data["refresh_token"])
            if access_token:
                return access_token
        
        # Try to login with test account
        if role in self.test_accounts:
            account = self.test_accounts[role]
            access_token = self._login(account["username"], account["password"], role)
            if access_token:
                return access_token
            
            # If login failed, try to register
            print(f"Login failed for {account['username']}, attempting to register...")
            if self._register(account["username"], account["password"], role):
                # Try login again after registration
                access_token = self._login(account["username"], account["password"], role)
                if access_token:
                    return access_token
        
        print(f"Failed to get access token for role: {role}")
        return None
    
    def call_api(self, endpoint: str, method: str = "GET", role: str = "USER", 
                 data: Optional[Dict] = None, params: Optional[Dict] = None) -> Tuple[int, Dict]:
        """Call API endpoint with authentication"""
        # Get access token
        access_token = self.get_access_token(role)
        if not access_token:
            return 401, {"error": "Failed to authenticate"}
        
        # Prepare headers
        headers = {
            "Authorization": f"Bearer {access_token}",
            "Content-Type": "application/json"
        }
        
        # Make request
        url = f"{self.base_url}{endpoint}"
        
        try:
            response = requests.request(
                method=method,
                url=url,
                headers=headers,
                json=data,
                params=params
            )
            
            # Try to parse JSON response
            try:
                response_data = response.json()
            except:
                response_data = {"text": response.text}
            
            return response.status_code, response_data
            
        except ConnectionError:
            return 0, {"error": f"Failed to connect to {url}"}
        except Exception as e:
            return 0, {"error": str(e)}


def main():
    parser = argparse.ArgumentParser(
        description="Resort Management System API Test Client",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Simple GET request
  %(prog)s /api/v1/users
  
  # POST request with admin role
  %(prog)s /api/v1/dev/generate-essential-data -m POST -r SUPER_ADMIN
  
  # POST with data
  %(prog)s /api/v1/reservations -m POST -d '{"roomId": 1, "guestName": "Test Guest"}'
  
  # GET with query parameters
  %(prog)s /api/v1/reservations -p page=0 -p size=20
  
  # Check API health
  %(prog)s /api/v1/health
        """
    )
    
    parser.add_argument("endpoint", help="API endpoint (e.g., /api/v1/users)")
    parser.add_argument("-m", "--method", default="GET", 
                       choices=["GET", "POST", "PUT", "PATCH", "DELETE"],
                       help="HTTP method (default: GET)")
    parser.add_argument("-r", "--role", default="USER",
                       choices=["USER", "ADMIN", "SUPER_ADMIN"],
                       help="User role for authentication (default: USER)")
    parser.add_argument("-d", "--data", type=str,
                       help="JSON data for request body")
    parser.add_argument("-p", "--param", action="append",
                       help="Query parameters (can be used multiple times, format: key=value)")
    parser.add_argument("-u", "--url", default="http://localhost:8080",
                       help="Base URL for API (default: http://localhost:8080)")
    parser.add_argument("-v", "--verbose", action="store_true",
                       help="Verbose output")
    
    args = parser.parse_args()
    
    # Parse query parameters
    params = {}
    if args.param:
        for param in args.param:
            if "=" in param:
                key, value = param.split("=", 1)
                params[key] = value
    
    # Parse JSON data
    data = None
    if args.data:
        try:
            data = json.loads(args.data)
        except json.JSONDecodeError:
            print(f"Error: Invalid JSON data: {args.data}")
            sys.exit(1)
    
    # Create client and make request
    client = APITestClient(base_url=args.url)
    
    if args.verbose:
        print(f"Calling {args.method} {args.endpoint}")
        if params:
            print(f"Parameters: {params}")
        if data:
            print(f"Data: {json.dumps(data, indent=2)}")
        print("-" * 50)
    
    status_code, response_data = client.call_api(
        endpoint=args.endpoint,
        method=args.method,
        role=args.role,
        data=data,
        params=params
    )
    
    # Print response
    if status_code == 0:
        print(f"Error: {response_data.get('error', 'Unknown error')}")
        sys.exit(1)
    else:
        if args.verbose:
            print(f"Status: {status_code}")
        print(json.dumps(response_data, indent=2, ensure_ascii=False))
        
        # Exit with non-zero code for HTTP errors
        if status_code >= 400:
            sys.exit(1)


if __name__ == "__main__":
    main()