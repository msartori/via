import axios from 'axios';

const apiUrl = import.meta.env.VITE_API_URL;

export async function getGuideToWidthraw(viaGuideId) {
  try {
    const res = await axios.get(`${apiUrl}/guide-to-withdraw/${viaGuideId}`, {
      headers: {
        'Content-Type': 'application/json',
        'bypass-tunnel-reminder': 'true',
        'Accept-Language': 'es',
      },
      validateStatus: () => true // will allow manage status
    });

    return {
      status: res.status,
      content: res.data
    };
  } catch (err) {
    console.error(err)
    return {
      status: 500,
      content: {
        message: 'Error al conectarse al servidor.',
        requestId: null
      }
    };
  }
}

export async function createGuideToWidthraw(viaGuideId) {
  try {
    const res = await axios.post(`${apiUrl}/guide-to-withraw`, {
      viaGuideId: viaGuideId 
    }, {
      headers: {
        'Content-Type': 'application/json',
        'bypass-tunnel-reminder': 'true',
        'Accept-Language': 'es',
      },
      validateStatus: () => true // will allow manage status
    });

    return {
      status: res.status,
      content: res.data
    };
  } catch (err) {
    console.error(err)
    return {
      status: 500,
      content: {
        message: 'Error al conectarse al servidor.',
        requestId: null
      }
    };
  }
}

export async function getMonitorEvents(params) {
  try {
    const res = await axios.get(`${apiUrl}/monitor/events`, {
      headers: {
        'Content-Type': 'application/json',
        'bypass-tunnel-reminder': 'true',
        'Accept-Language': 'es',
      },
      validateStatus: () => true // will allow manage status
    });
    return {
      status: res.status,
      content: res.data
    };
  } catch (err) {
    console.error(err)
    return {
      status: 500,
      content: {
        message: 'Error al conectarse al servidor.',
        requestId: null
      }
    };
  }
}

export async function getOperatorGuides(operatorId) {
  try {
    const res = await axios.get(`${apiUrl}/operator/guides`, { withCredentials: true }, {
      params: {
        operatorId
      },
      headers: {
        'Content-Type': 'application/json',
        'bypass-tunnel-reminder': 'true',
        'Accept-Language': 'es',
      },
      validateStatus: () => true // permite manejar manualmente el status
    });

    return {
      status: res.status,
      content: res.data
    };
  } catch (err) {
    console.error(err);
    return {
      status: 500,
      content: {
        message: 'Error al conectarse al servidor.',
        requestId: null
      }
    };
  }
}

export async function assignGuideToOperator(guideId, operatorId) {
  try {
    const res = await axios.post(`${apiUrl}/guide/${guideId}/assign`, {
      operatorId: Number(operatorId)
    }, {
      headers: {
        'Content-Type': 'application/json',
        'bypass-tunnel-reminder': 'true',
        'Accept-Language': 'es',
      },
      validateStatus: () => true
    });

    return {
      status: res.status,
      content: res.data
    };
  } catch (err) {
    console.error(err);
    return {
      status: 500,
      content: {
        message: 'Error al asignar la guía al operador.',
        requestId: null
      }
    };
  }
}

export async function getGuideStatusOptions(guideId) {
  try {
    const res = await axios.get(`${apiUrl}/guide/${guideId}/status-options`, {
      headers: {
        'Content-Type': 'application/json',
        'bypass-tunnel-reminder': 'true',
        'Accept-Language': 'es',
      },
      validateStatus: () => true
    });

    return {
      status: res.status,
      content: res.data
    };
  } catch (err) {
    console.error(err)
    return {
      status: 500,
      content: {
        message: 'Error al obtener opciones de estado',
        requestId: null
      }
    };
  }
}

export async function changeGuideStatus(guideId, newStatusId) {
  try {
    const response = await fetch(`${apiUrl}/guide/${guideId}/status`, {
      method: 'PUT',
      headers: { 
        'Content-Type': 'application/json',
        'bypass-tunnel-reminder': 'true',
        'Accept-Language': 'es',
      },
      body: JSON.stringify({ status: newStatusId })
    })
    const content = await response.json()
    return { status: response.status, content }
  } catch (err) {
    console.error(err)
    return {
      status: 500,
      content: {
        message: 'Error al obtener opciones de estado',
        requestId: null
      }
    };
  }
}
