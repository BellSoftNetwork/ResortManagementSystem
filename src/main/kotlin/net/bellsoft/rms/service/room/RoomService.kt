package net.bellsoft.rms.service.room

import mu.KLogging
import net.bellsoft.rms.component.auth.SecuritySupport
import net.bellsoft.rms.component.history.EntityHistoryComponent
import net.bellsoft.rms.component.history.dto.EntityHistoryDto
import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.domain.room.RoomRepository
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.DuplicateDataException
import net.bellsoft.rms.exception.UserNotFoundException
import net.bellsoft.rms.mapper.model.RoomMapper
import net.bellsoft.rms.service.common.dto.EntityListDto
import net.bellsoft.rms.service.room.dto.RoomCreateDto
import net.bellsoft.rms.service.room.dto.RoomDetailDto
import net.bellsoft.rms.service.room.dto.RoomFilterDto
import net.bellsoft.rms.service.room.dto.RoomPatchDto
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class RoomService(
    private val roomRepository: RoomRepository,
    private val entityHistoryComponent: EntityHistoryComponent,
    private val securitySupport: SecuritySupport,
    private val roomMapper: RoomMapper,
) {
    fun findAll(pageable: Pageable, filter: RoomFilterDto) = EntityListDto
        .of(
            roomRepository.getFilteredRooms(pageable, filter),
            roomMapper::toDto,
        )

    fun find(id: Long) = roomMapper.toDto(
        roomRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 객실"),
    )

    @Transactional
    fun create(request: RoomCreateDto): RoomDetailDto {
        try {
            return roomMapper.toDto(roomRepository.save(roomMapper.toEntity(request)))
        } catch (e: DataIntegrityViolationException) {
            logger.warn(e.message)
            throw DuplicateDataException("이미 존재하는 객실")
        }
    }

    @Transactional
    fun update(id: Long, request: RoomPatchDto): RoomDetailDto {
        val room = roomRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 객실")

        roomMapper.updateEntity(request, room)

        return roomMapper.toDto(roomRepository.save(room))
    }

    @Transactional
    fun delete(id: Long) {
        val room = roomRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 객실")
        val user = securitySupport.getCurrentUser()
            ?: throw UserNotFoundException("로그인 필요")

        roomRepository.save(room.apply { updatedBy = user })
        roomRepository.flush()
        roomRepository.delete(room)
    }

    fun findHistory(id: Long, pageable: Pageable): EntityListDto<EntityHistoryDto<RoomDetailDto>> = EntityListDto
        .of(entityHistoryComponent.findAllHistory<Room, RoomDetailDto>(Room::class, roomMapper::toDto, id, pageable))

    companion object : KLogging()
}
