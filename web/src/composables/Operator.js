import { ref, nextTick, onMounted, onBeforeUnmount, computed, handleError } from 'vue'
import {assignGuideToOperator, getGuideStatusOptions, changeGuideStatus, handleAuthRedirect, doLogout } from '../services/api'
import {apiSSEUrl, apiSSE} from '../services/apiConfig'

export default function useOperator({
  activityPanel,
  showConfirmModal,
  showSuccessModal,
  error,
  requestId,
  statusChanging,
  animateActivity,
  loggingOut
}) {
  const operatorGuides = ref([])
  const activeGuide = ref(null)
  const statusOptions = ref([])
  const loadingStatusOptions = ref(false)
  const pendingStatusChange = ref({ guideId: null, status: null, viaGuideId: null })
  

  let operatorGuidesSource = null

  const verifyUnauthorizedOnSSE = async (uri) => {
    try {
      const res = await apiSSE.get(uri, {
        headers: {
          Accept: 'text/event-stream',
        },
        validateStatus: () => true
      })
      return res.status
    } catch (e) {
      console.error('SSE pre-check failed', e)
      return 500
    }
  }

  const fetchOperatorGuides = async () => {
    operatorGuidesSource = new EventSource(`${apiSSEUrl}/operator/guides?lang=es`,  { withCredentials: true })
    operatorGuidesSource.onmessage = (event) => {
      try {
        const content = JSON.parse(event.data)
        const updatedGuides = content.data?.operatorGuides || []
        operatorGuides.value = updatedGuides
        if (activeGuide.value) {
          const stillExists = updatedGuides.find(g => g.guideId === activeGuide.value.guideId)
          if (stillExists && stillExists.operator.id != 0) {
            activeGuide.value = { ...stillExists }
            loadStatusOptions(stillExists.guideId)
          } else {
            activeGuide.value = null
          }
        } 
        if (content.message != "") {
          //console.error(content.message)
          error.value = content.message
          requestId.value = content.requestId
        }
      } catch (e) {
        console.error('Failed to parse SSE data', e)
        handleError('Fallo al interpretar data de SSE')
      }
    }

    operatorGuidesSource.onerror = async (err) => {
      let state = await verifyUnauthorizedOnSSE('/operator/guides?lang=es')
      handleAuthRedirect(state)
      //console.error('SSE connection error:', err)
      //error.value = 'Error al conectarse al servidor'
      //requestId.value = ''
      operatorGuidesSource.close()
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
    if (guide.operator.id != 1) {
      activeGuide.value = { ...guide }
      await loadStatusOptions(guide.guideId)
      scrollToActivity()
      return
    }

    const response = await assignGuideToOperator(guide.guideId)
    if (response.status === 200) {
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

  async function logout() {
    loggingOut.value = true
    const res = await doLogout()
    if (res.status === 200) {
      if (res.content.message == "") {
        window.location.reload()
      } else {
          error.value = res.content.message
          requestId.value = res.content.requestId
      }
    }
    loggingOut.value = false
  }

  onMounted(() => {
    fetchOperatorGuides()
  })

  onBeforeUnmount(() => {
    operatorGuidesSource.close()
  })

  setInterval(() => {
    if (activeGuide.value?.lastChange) activeGuide.value = { ...activeGuide.value }
  }, 1000)

  return {
    operatorGuides,
    activeGuide,
    selectGuide,
    openConfirmModal,
    closeConfirmModal,
    confirmChangeStatus,
    closeSuccessModal,
    statusOptions,
    loadingStatusOptions,
    pendingStatusChange,
    elapsedTime,
    logout
  }
}
