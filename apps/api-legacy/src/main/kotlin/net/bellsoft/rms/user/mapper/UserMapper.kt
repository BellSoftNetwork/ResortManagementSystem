package net.bellsoft.rms.user.mapper

import net.bellsoft.rms.common.config.BaseMapperConfig
import net.bellsoft.rms.common.util.MD5Util
import net.bellsoft.rms.user.dto.response.UserDetailDto
import net.bellsoft.rms.user.dto.response.UserSummaryDto
import net.bellsoft.rms.user.entity.User
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.Mappings
import org.mapstruct.Named

@Mapper(config = BaseMapperConfig::class)
abstract class UserMapper {
    @Mappings(
        Mapping(target = "profileImageUrl", source = "email", qualifiedByName = ["emailToProfileImageUrl"]),
    )
    abstract fun toDto(entity: User): UserDetailDto

    abstract fun toDto(dto: UserDetailDto): UserSummaryDto

    @Named("emailToProfileImageUrl")
    fun emailToProfileImageUrl(email: String?) =
        email?.let { "https://gravatar.com/avatar/${MD5Util.md5Hex(email)}" }
            ?: "https://gravatar.com/avatar/00000000000000000000"
}
