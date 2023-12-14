package net.bellsoft.rms.room.controller

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import net.bellsoft.rms.common.dto.response.ListResponse
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.room.dto.filter.RoomRequestFilter
import net.bellsoft.rms.room.dto.request.RoomGroupCreateRequest
import net.bellsoft.rms.room.dto.request.RoomGroupPatchRequest
import net.bellsoft.rms.room.dto.response.RoomGroupDetailDto
import net.bellsoft.rms.room.dto.response.RoomGroupSummaryDto
import net.bellsoft.rms.user.entity.User
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.security.access.annotation.Secured
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.DeleteMapping
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PatchMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.ResponseStatus
import org.springframework.web.bind.annotation.RestController

@Tag(name = "객실 그룹", description = "객실 그룹 관리 API")
@SecurityRequirement(name = "basicAuth")
@Validated
@RestController
@RequestMapping("/api/v1/room-groups")
interface RoomGroupController {
    @Operation(summary = "객실 그룹 리스트", description = "객실 그룹 리스트 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getRoomGroups(pageable: Pageable): ResponseEntity<ListResponse<RoomGroupSummaryDto>>

    @Operation(summary = "객실 그룹 조회", description = "객실 그룹 단건 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/{id}")
    fun getRoomGroup(
        @PathVariable("id") id: Long,
        filter: RoomRequestFilter,
    ): ResponseEntity<SingleResponse<RoomGroupDetailDto>>

    @Operation(summary = "객실 그룹 생성", description = "객실 그룹 생성")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @PostMapping
    fun createRoomGroup(
        @AuthenticationPrincipal user: User,

        @RequestBody @Valid
        request: RoomGroupCreateRequest,
    ): SingleResponse<RoomGroupSummaryDto>

    @Operation(summary = "객실 그룹 수정", description = "기존 객실 그룹 정보 수정")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @PatchMapping("/{id}")
    fun updateRoomGroup(
        @AuthenticationPrincipal user: User,
        @PathVariable("id") id: Long,

        @RequestBody @Valid
        request: RoomGroupPatchRequest,
    ): ResponseEntity<SingleResponse<RoomGroupSummaryDto>>

    @Operation(summary = "객실 그룹 삭제", description = "기존 객실 그룹 삭제")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "204"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    fun deleteRoomGroup(
        @AuthenticationPrincipal user: User,
        @PathVariable("id") id: Long,
    )
}
