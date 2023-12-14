package net.bellsoft.rms.room.service.impl

import mu.KLogging
import net.bellsoft.rms.authentication.component.SecuritySupport
import net.bellsoft.rms.authentication.exception.UserNotFoundException
import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.exception.DuplicateDataException
import net.bellsoft.rms.common.exception.RelatedDataException
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.response.RoomGroupDetailDto
import net.bellsoft.rms.room.dto.response.RoomGroupSummaryDto
import net.bellsoft.rms.room.dto.service.RoomGroupCreateDto
import net.bellsoft.rms.room.dto.service.RoomGroupPatchDto
import net.bellsoft.rms.room.mapper.RoomGroupMapper
import net.bellsoft.rms.room.mapper.RoomProjectionMapper
import net.bellsoft.rms.room.repository.RoomGroupRepository
import net.bellsoft.rms.room.repository.RoomRepository
import net.bellsoft.rms.room.service.RoomGroupService
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.data.repository.findByIdOrNull
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class RoomGroupServiceImpl(
    private val securitySupport: SecuritySupport,
    private val roomGroupRepository: RoomGroupRepository,
    private val roomRepository: RoomRepository,
    private val roomGroupMapper: RoomGroupMapper,
    private val roomProjectionMapper: RoomProjectionMapper,
) : RoomGroupService {
    override fun findAll(pageable: Pageable): EntityListDto<RoomGroupSummaryDto> = EntityListDto
        .of(
            roomGroupRepository.findAll(pageable),
            roomGroupMapper::toDto,
        )

    override fun find(id: Long, filter: RoomFilterDto): RoomGroupDetailDto {
        val roomGroup = getRoomGroup(id)
        val roomLastStayProjections = roomGroupRepository.getFilteredRoomsOrderByLastStayAt(roomGroup, filter)

        return roomProjectionMapper.toDto(getRoomGroup(id), roomLastStayProjections)
    }

    @Transactional
    override fun create(request: RoomGroupCreateDto): RoomGroupSummaryDto {
        try {
            return roomGroupMapper.toDto(roomGroupRepository.save(roomGroupMapper.toEntity(request)))
        } catch (e: DataIntegrityViolationException) {
            logger.warn(e.message)
            throw DuplicateDataException("이미 존재하는 객실 그룹")
        }
    }

    @Transactional
    override fun update(id: Long, request: RoomGroupPatchDto): RoomGroupSummaryDto {
        val roomGroup = getRoomGroup(id)

        roomGroupMapper.updateEntity(request, roomGroup)

        return roomGroupMapper.toDto(roomGroupRepository.save(roomGroup))
    }

    @Transactional
    override fun delete(id: Long) {
        val user = securitySupport.getCurrentUser()
            ?: throw UserNotFoundException("로그인 필요")
        val roomGroup = getRoomGroup(id)

        if (roomRepository.existsByRoomGroup(roomGroup))
            throw RelatedDataException("그룹 내 객실이 존재하여 삭제 불가")

        roomGroupRepository.save(roomGroup.apply { updatedBy = user })
        roomGroupRepository.flush()
        roomGroupRepository.delete(roomGroup)
    }

    private fun getRoomGroup(id: Long) = roomGroupRepository.findByIdOrNull(id)
        ?: throw DataNotFoundException("존재하지 않는 객실 그룹")

    companion object : KLogging()
}
