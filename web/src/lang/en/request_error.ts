const requestError = {
  networkError: 'Backend API error',
  timeout: 'System API request timeout',
  requestFailed: 'API error',
  default: 'Unknown error',
  400: 'Invalid request',
  401: 'No permission. Please log in again',
  403: 'Access denied',
  404: "Request error. Can't find requested resources",
  405: 'Request method not allowed',
  408: 'Request timeout',
  500: 'Server error',
  501: 'Network not implemented',
  502: 'Network error',
  503: 'Service unavailable',
  504: 'Network timeout',
  505: 'Request not supported by HTTP',
  other: 'Other connection error --{{code}}',
};

export default requestError;
