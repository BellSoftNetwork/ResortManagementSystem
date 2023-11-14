package net.bellsoft.rms.service.room

import mu.KLogging
import net.bellsoft.rms.component.auth.SecuritySupport
import net.bellsoft.rms.component.history.EntityHistoryComponent
import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.domain.room.RoomRepository
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.DuplicateDataException
import net.bellsoft.rms.exception.UserNotFoundException
import net.bellsoft.rms.service.common.dto.EntityListDto
import net.bellsoft.rms.service.room.dto.RoomCreateDto
import net.bellsoft.rms.service.room.dto.RoomDto
import net.bellsoft.rms.service.room.dto.RoomUpdateDto
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
class RoomService(
    private val roomRepository: RoomRepository,
    private val entityHistoryComponent: EntityHistoryComponent,
    private val securitySupport: SecuritySupport,
) {
    fun findAll(pageable: Pageable) = EntityListDto
        .of(
            roomRepository.findAll(pageable),
            RoomDto::of,
        )

    fun find(id: Long) = RoomDto.of(
        roomRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 객실"),
    )

    @Transactional
    fun create(request: RoomCreateDto): RoomDto {
        try {
            return RoomDto.of(roomRepository.save(request.toEntity()))
        } catch (e: DataIntegrityViolationException) {
            logger.warn(e.message)
            throw DuplicateDataException("이미 존재하는 객실")
        }
    }

    @Transactional
    fun update(id: Long, request: RoomUpdateDto): RoomDto {
        val room = roomRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 객실")

        request.updateEntity(room)

        return RoomDto.of(roomRepository.save(room))
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

    fun findHistory(id: Long, pageable: Pageable) = EntityListDto
        .of(entityHistoryComponent.findAllHistory(Room::class, RoomDto::of, id, pageable))

    companion object : KLogging()
}
