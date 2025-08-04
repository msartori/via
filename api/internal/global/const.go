package global

type contextKey string

const RequestIDKey contextKey = "requestID"

// Pub Sub Channels
const NewGuideChannel string = "new_guide"
const GuideStatusChangeChannel string = "guide_status_change"
const GuideAssignmentChannel string = "guide_assignment"
