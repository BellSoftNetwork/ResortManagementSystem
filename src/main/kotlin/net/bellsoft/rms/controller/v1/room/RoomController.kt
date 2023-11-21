package net.bellsoft.rms.controller.v1.room

import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.responses.ApiResponses
import io.swagger.v3.oas.annotations.security.SecurityRequirement
import io.swagger.v3.oas.annotations.tags.Tag
import jakarta.validation.Valid
import net.bellsoft.rms.controller.common.dto.ListResponse
import net.bellsoft.rms.controller.common.dto.SingleResponse
import net.bellsoft.rms.controller.v1.room.dto.RoomCreateRequest
import net.bellsoft.rms.controller.v1.room.dto.RoomRequestFilter
import net.bellsoft.rms.controller.v1.room.dto.RoomUpdateRequest
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.service.room.RoomService
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
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
class RoomController(
    private val roomService: RoomService,
) {
    @Operation(summary = "객실 리스트", description = "객실 리스트 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping
    fun getRooms(pageable: Pageable, filter: RoomRequestFilter) = ListResponse
        .of(roomService.findAll(pageable, filter.toDto()), filter)
        .toResponseEntity()

    @Operation(summary = "객실 조회", description = "객실 단건 조회")
    @ApiResponses(
        value = [
            ApiResponse(responseCode = "200"),
        ],
    )
    @GetMapping("/{id}")
    fun getRoom(@PathVariable("id") id: Long) = SingleResponse
        .of(roomService.find(id))
        .toResponseEntity()

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
    ) = SingleResponse
        .of(roomService.create(request.toDto()))

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
        request: RoomUpdateRequest,
    ) = SingleResponse
        .of(roomService.update(id, request.toDto()))
        .toResponseEntity()

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
    ) {
        roomService.delete(id)
    }

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
    ) = ListResponse
        .of(roomService.findHistory(id, pageable))
        .toResponseEntity()
}
