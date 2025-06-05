package net.bellsoft.rms.room.controller

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import net.bellsoft.rms.common.dto.response.ListResponse
import net.bellsoft.rms.common.dto.response.SingleResponse
import net.bellsoft.rms.revision.dto.EntityRevisionDto
import net.bellsoft.rms.room.dto.filter.RoomRequestFilter
import net.bellsoft.rms.room.dto.request.RoomCreateRequest
import net.bellsoft.rms.room.dto.request.RoomPatchRequest
import net.bellsoft.rms.room.dto.response.RoomDetailDto
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

@Tag(name = "객실", description = "객실 관리 API")
@SecurityRequirement(name = "basicAuth")
@Validated
@RestController
@RequestMapping("/api/v1/rooms")
interface RoomController {
    @Operation(summary = "객실 리스트", description = "객실 리스트 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getRooms(pageable: Pageable, filter: RoomRequestFilter): ResponseEntity<ListResponse<RoomDetailDto>>

    @Operation(summary = "객실 조회", description = "객실 단건 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/{id}")
    fun getRoom(@PathVariable("id") id: Long): ResponseEntity<SingleResponse<RoomDetailDto>>

    @Operation(summary = "객실 생성", description = "객실 생성")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "201"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @PostMapping
    fun createRoom(
        @AuthenticationPrincipal user: User,

        @RequestBody @Valid
        request: RoomCreateRequest,
    ): SingleResponse<RoomDetailDto>

    @Operation(summary = "객실 수정", description = "기존 객실 정보 수정")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @PatchMapping("/{id}")
    fun updateRoom(
        @AuthenticationPrincipal user: User,
        @PathVariable("id") id: Long,

        @RequestBody @Valid
        request: RoomPatchRequest,
    ): ResponseEntity<SingleResponse<RoomDetailDto>>

    @Operation(summary = "객실 삭제", description = "기존 객실 삭제")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "204"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    fun deleteRoom(
        @AuthenticationPrincipal user: User,
        @PathVariable("id") id: Long,
    )

    @Operation(summary = "객실 이력", description = "객실 정보 변경 이력 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @Secured("ADMIN", "SUPER_ADMIN")
    @GetMapping("/{id}/histories")
    fun getRoomHistory(
        @PathVariable("id") id: Long,
        pageable: Pageable,
    ): ResponseEntity<ListResponse<EntityRevisionDto<RoomDetailDto>>>
}
