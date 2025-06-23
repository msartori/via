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
    const res = await axios.get(`${apiUrl}/operator/guides`, {
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