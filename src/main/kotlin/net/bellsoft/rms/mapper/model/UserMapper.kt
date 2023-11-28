package net.bellsoft.rms.mapper.model

import net.bellsoft.rms.controller.v1.admin.dto.AdminUserCreateRequest
import net.bellsoft.rms.controller.v1.auth.dto.UserRegistrationRequest
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.mapper.common.JsonNullableMapper
import net.bellsoft.rms.mapper.common.ReferenceMapper
import net.bellsoft.rms.service.auth.dto.UserCreateDto
import net.bellsoft.rms.service.auth.dto.UserDetailDto
import net.bellsoft.rms.util.MD5Util
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.Mappings
import org.mapstruct.Named
import org.mapstruct.NullValuePropertyMappingStrategy

@Mapper(
    uses = [JsonNullableMapper::class, ReferenceMapper::class],
    nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE,
    componentModel = "spring",
)
abstract class UserMapper {
    @Mappings(
        Mapping(target = "profileImageUrl", source = "email", qualifiedByName = ["emailToProfileImageUrl"]),
    )
    abstract fun toDto(entity: User): UserDetailDto

    abstract fun toDto(dto: UserRegistrationRequest): UserCreateDto
    abstract fun toDto(dto: AdminUserCreateRequest): UserCreateDto

    @Named("emailToProfileImageUrl")
    fun emailToProfileImageUrl(email: String) = "https://gravatar.com/avatar/${MD5Util.md5Hex(email)}"
}
