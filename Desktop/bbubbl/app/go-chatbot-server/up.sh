# 패스워드 변수 설정 (보안을 위해 환경변수에서 가져오는 것을 권장)
SSH_PASS="your_password_here"

log "프로젝트 파일 복사 시작"
sshpass -p "$SSH_PASS" rsync -avz ./ $SSH_HOST:$REMOTE_DIR --exclude qdrant_storage
log "프로젝트 파일 복사 완료"

log "AWS 서버에 SSH 접속 시작"
sshpass -p "$SSH_PASS" ssh $SSH_HOST << EOF
// ... existing code ... 