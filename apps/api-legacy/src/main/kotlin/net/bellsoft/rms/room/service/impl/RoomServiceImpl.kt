package net.bellsoft.rms.room.service.impl

import mu.KLogging
import net.bellsoft.rms.authentication.component.SecuritySupport
import net.bellsoft.rms.authentication.exception.UserNotFoundException
import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.exception.DuplicateDataException
import net.bellsoft.rms.revision.component.EntityRevisionComponent
import net.bellsoft.rms.revision.dto.EntityRevisionDto
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.response.RoomDetailDto
import net.bellsoft.rms.room.dto.service.RoomCreateDto
import net.bellsoft.rms.room.dto.service.RoomPatchDto
import net.bellsoft.rms.room.entity.Room
import net.bellsoft.rms.room.mapper.RoomMapper
import net.bellsoft.rms.room.repository.RoomRepository
import net.bellsoft.rms.room.service.RoomService
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class RoomServiceImpl(
    private val roomRepository: RoomRepository,
    private val entityRevisionComponent: EntityRevisionComponent,
    private val securitySupport: SecuritySupport,
    private val roomMapper: RoomMapper,
) : RoomService {
    override fun findAll(pageable: Pageable, filter: RoomFilterDto) = EntityListDto
        .of(
            roomRepository.getFilteredRooms(pageable, filter),
            roomMapper::toDto,
        )

    override fun find(id: Long) = roomMapper.toDto(
        roomRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 객실"),
    )

    @Transactional
    override fun create(request: RoomCreateDto): RoomDetailDto {
        try {
            return roomMapper.toDto(roomRepository.save(roomMapper.toEntity(request)))
        } catch (e: DataIntegrityViolationException) {
            logger.warn(e.message)
            throw DuplicateDataException("이미 존재하는 객실")
        }
    }

    @Transactional
    override fun update(id: Long, request: RoomPatchDto): RoomDetailDto {
        val room = roomRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 객실")

        roomMapper.updateEntity(request, room)

        return roomMapper.toDto(roomRepository.save(room))
    }

    @Transactional
    override fun delete(id: Long) {
        val room = roomRepository.findByIdOrNull(id)
            ?: throw DataNotFoundException("존재하지 않는 객실")
        val user = securitySupport.getCurrentUser()
            ?: throw UserNotFoundException("로그인 필요")

        roomRepository.save(room.apply { updatedBy = user })
        roomRepository.flush()
        roomRepository.delete(room)
    }

    override fun findHistory(id: Long, pageable: Pageable): EntityListDto<EntityRevisionDto<RoomDetailDto>> =
        EntityListDto
            .of(
                entityRevisionComponent.findAllHistory<Room, RoomDetailDto>(
                    Room::class,
                    roomMapper::toDto,
                    id,
                    pageable,
                ),
            )

    companion object : KLogging()
}
