---
id: api-legacy-user-management
title: "api-legacy 사용자 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: small
---

# api-legacy 사용자 관리

> 사용자 프로필 및 관리자 계정 관리

---

## 1. 엔드포인트

### 1.1 프로필

| Method | Path | 설명 |
|--------|------|------|
| GET | `/api/v1/my` | 내 정보 조회 |
| POST | `/api/v1/my` | 내 정보 조회 (호환용) |
| PATCH | `/api/v1/my` | 내 정보 수정 |

### 1.2 관리자 계정

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/admin/accounts` | 사용자 목록 | ADMIN |
| POST | `/api/v1/admin/accounts` | 사용자 생성 | ADMIN |
| PATCH | `/api/v1/admin/accounts/{id}` | 사용자 수정 | ADMIN |

---

## 2. User 엔티티

```kotlin
@Entity
@Table(
    uniqueConstraints = [
        UniqueConstraint(columnNames = ["user_id", "deleted_at"]),
        UniqueConstraint(columnNames = ["email", "deleted_at"])
    ]
)
@SQLDelete(sql = "UPDATE user SET deleted_at = NOW() WHERE id = ?")
@Where(clause = "deleted_at = '1970-01-01 00:00:00'")
class User(
    @Column(name = "user_id", nullable = false, length = 30)
    var userId: String,
    
    @Column(nullable = true, length = 100)
    var email: String? = null,
    
    @Column(nullable = false, length = 20)
    var name: String,
    
    @Column(nullable = false, length = 100)
    private var password: String,
    
    @Column(nullable = false, columnDefinition = "TINYINT")
    var status: UserStatus = UserStatus.INACTIVE,
    
    @Column(nullable = false, columnDefinition = "TINYINT")
    var role: UserRole = UserRole.NORMAL
) : BaseTimeEntity(), UserDetails
```

---

## 3. 역할 및 상태

### 3.1 UserRole

```kotlin
enum class UserRole(val value: Int) {
    NORMAL(0),
    ADMIN(100),
    SUPER_ADMIN(127)
}
```

### 3.2 UserStatus

```kotlin
enum class UserStatus(val value: Int) {
    INACTIVE(-1),
    ACTIVE(1)
}
```

---

## 4. 서비스

### 4.1 UserService 인터페이스

```kotlin
interface UserService {
    fun register(request: UserRegistrationRequest): UserDetailDto
    fun getAll(pageable: Pageable): EntityListDto<User, UserDetailDto>
    fun getById(id: Long): UserDetailDto
    fun create(request: AdminUserCreateRequest): UserDetailDto
    fun update(id: Long, request: AdminUserPatchRequest): UserDetailDto
    fun updateMy(request: MyPatchRequest): UserDetailDto
}
```

### 4.2 권한 체크

```kotlin
@PreAuthorize("hasAnyRole('ADMIN', 'SUPER_ADMIN')")
fun getAll(pageable: Pageable): ListResponse<UserDetailDto>

@PreAuthorize("hasRole('SUPER_ADMIN') or @securitySupport.isCurrentUser(#id)")
fun update(id: Long, request: AdminUserPatchRequest): SingleResponse<UserDetailDto>
```

---

## 5. DTO

### 5.1 요청

```kotlin
data class AdminUserCreateRequest(
    val userId: String,
    val email: String?,
    val name: String,
    val password: String,
    val role: UserRole?,
    val status: UserStatus?
)

data class AdminUserPatchRequest(
    val email: JsonNullable<String?>,
    val name: JsonNullable<String>?,
    val status: JsonNullable<UserStatus>?,
    val role: JsonNullable<UserRole>?
)

data class MyPatchRequest(
    val name: JsonNullable<String>?,
    val currentPassword: String?,
    val newPassword: String?
)
```

### 5.2 응답

```kotlin
data class UserDetailDto(
    val id: Long,
    val userId: String,
    val email: String?,
    val name: String,
    val status: UserStatus,
    val role: UserRole,
    val profileImageUrl: String,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime
)
```

---

## 6. MapStruct 매퍼

```kotlin
@Mapper(componentModel = "spring")
interface UserMapper {
    fun toDetailDto(entity: User): UserDetailDto
    
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "password", ignore = true)
    fun toEntity(request: AdminUserCreateRequest): User
    
    @BeanMapping(nullValuePropertyMappingStrategy = IGNORE)
    fun updateEntity(@MappingTarget entity: User, request: AdminUserPatchRequest)
}
```

---

## 7. 비밀번호 변경

```kotlin
// MyPatchRequest 처리
fun updateMy(request: MyPatchRequest): UserDetailDto {
    val user = getCurrentUser()
    
    if (request.newPassword != null) {
        if (request.currentPassword == null) {
            throw BadRequestException("현재 비밀번호를 입력해주세요")
        }
        if (!passwordEncoder.matches(request.currentPassword, user.password)) {
            throw BadRequestException("현재 비밀번호가 일치하지 않습니다")
        }
        user.password = passwordEncoder.encode(request.newPassword)
    }
    
    // 다른 필드 업데이트
    return userMapper.toDetailDto(userRepository.save(user))
}
```
