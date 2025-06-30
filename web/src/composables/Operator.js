import { ref, nextTick, onMounted, onBeforeUnmount, computed } from 'vue'
import { getOperatorGuides, assignGuideToOperator, getGuideStatusOptions, changeGuideStatus } from '../services/api'

export default function useOperator({
  activityPanel,
  showConfirmModal,
  showSuccessModal,
  error,
  requestId,
  statusChanging,
  animateActivity
}) {
  const operatorId = '3'
  const operatorGuides = ref([])
  const activeGuide = ref(null)
  const statusOptions = ref([])
  const loadingStatusOptions = ref(false)
  const pendingStatusChange = ref({ guideId: null, status: null, viaGuideId: null })

  let refreshIntervalId = null

  async function fetchOperatorGuides() {
    error.value = null
    const response = await getOperatorGuides(operatorId)
    const { status, content } = response

    if (status === 200) {
      const updatedGuides = content.data?.operatorGuides || []
      operatorGuides.value = updatedGuides

      if (activeGuide.value) {
        const stillExists = updatedGuides.find(g => g.guideId === activeGuide.value.guideId)
        if (stillExists && stillExists.operator?.id === Number(operatorId)) {
          activeGuide.value = { ...stillExists }
          await loadStatusOptions(stillExists.guideId)
        } else {
          activeGuide.value = null
        }
      }
    } else {
      error.value = content.message || 'Error al obtener las guías'
      requestId.value = content.requestId
    }
  }

  async function loadStatusOptions(guideId) {
    loadingStatusOptions.value = true
    statusOptions.value = []
    const response = await getGuideStatusOptions(guideId)
    const { status, content } = response
    if (status === 200) {
      statusOptions.value = content.data?.statusOption || []
    }
    loadingStatusOptions.value = false
  }

  function scrollToActivity() {
    nextTick(() => {
      activityPanel.value?.scrollIntoView({
        behavior: 'smooth',
        block: 'start',
      })
      if (animateActivity) {
        animateActivity.value = true
        setTimeout(() => (animateActivity.value = false), 800)
      }
    })
  }

  async function selectGuide(guide) {
    if (!guide.selectable) return

    if (guide.operator?.id === Number(operatorId)) {
      activeGuide.value = { ...guide }
      await loadStatusOptions(guide.guideId)
      scrollToActivity()
      return
    }

    const response = await assignGuideToOperator(guide.guideId, operatorId)
    if (response.status === 200) {
      await fetchOperatorGuides()
      const updated = operatorGuides.value.find(g => g.guideId === guide.guideId)
      if (updated) {
        activeGuide.value = { ...updated }
        await loadStatusOptions(updated.guideId)
        scrollToActivity()
      }
    }
  }

  function openConfirmModal(guideId, status, viaGuideId) {
    pendingStatusChange.value = { guideId, status, viaGuideId }
    showConfirmModal.value = true
  }

  function closeConfirmModal() {
    showConfirmModal.value = false
    pendingStatusChange.value = { guideId: null, status: null }
  }

  async function confirmChangeStatus() {
    const { guideId, status } = pendingStatusChange.value
    if (!guideId || !status?.id) return

    statusChanging.value = true
    const response = await changeGuideStatus(guideId, status.id)

    if (response.status === 200) {
      closeConfirmModal()
      showSuccessModal.value = true
    } else {
      alert('Error al cambiar el estado.')
    }

    statusChanging.value = false
  }

  async function closeSuccessModal() {
    showSuccessModal.value = false
    await fetchOperatorGuides()
  }

  const elapsedTime = computed(() => {
    if (!activeGuide.value?.lastChange) return ''
    const lastChange = new Date(activeGuide.value.lastChange)
    const now = new Date()
    const diff = now - lastChange
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    const hours = new Date(diff).getUTCHours().toString().padStart(2, '0')
    const minutes = new Date(diff).getUTCMinutes().toString().padStart(2, '0')
    const seconds = new Date(diff).getUTCSeconds().toString().padStart(2, '0')
    return `${days}d ${hours}:${minutes}:${seconds}`
  })

  onMounted(() => {
    fetchOperatorGuides()
    refreshIntervalId = setInterval(fetchOperatorGuides, 10000)
  })

  onBeforeUnmount(() => {
    clearInterval(refreshIntervalId)
  })

  setInterval(() => {
    if (activeGuide.value?.lastChange) activeGuide.value = { ...activeGuide.value }
  }, 1000)

  return {
    operatorGuides,
    activeGuide,
    selectGuide,
    fetchOperatorGuides,
    openConfirmModal,
    closeConfirmModal,
    confirmChangeStatus,
    closeSuccessModal,
    statusOptions,
    loadingStatusOptions,
    pendingStatusChange,
    elapsedTime
  }
}
