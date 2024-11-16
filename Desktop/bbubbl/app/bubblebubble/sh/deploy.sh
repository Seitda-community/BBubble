#!/bin/bash

# 로그 파일 경로 설정
LOG_FILE="./deployment.log"

# 로그 함수 정의: 콘솔과 파일에 동시에 로그 출력
log() {
    local message="[$(date '+%Y-%m-%d %H:%M:%S')] $1"
    echo "$message"
    echo "$message" >> "$LOG_FILE"
}

START_TIME=$(date +%s.%N)
log "스크립트 시작"

SSH_HOST="fye"
REMOTE_DIR="/home/ec2-user/next-meister"
PORT="3000"
log "변수 설정 완료"

# 현재 디렉토리 저장
CURRENT_DIR=$(pwd)
log "현재 디렉토리: $CURRENT_DIR"

cd $CURRENT_DIR

log "프로젝트 파일 복사 시작"
rsync -avz ./ $SSH_HOST:$REMOTE_DIR
log "프로젝트 파일 복사 완료"
exit 0
pnpm store prune && pnpm install && pnpm build

log "AWS 서버에 SSH 접속 시작"
ssh $SSH_HOST << EOF
    cd $REMOTE_DIR

    log() {
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] \$1"
    }
    
    log "원격 서버 작업 시작"
    
    # 3000 포트 사용 중인 프로세스 종료
    sudo kill -9 $(sudo lsof -t -i:3000)
    log "3000 포트 사용 중인 프로세스 종료 완료"

    log "Next.js 애플리케이션 백그라운드에서 시작 중"
    sudo nohup pnpm start > app.log 2>&1 &
    log "Next.js 애플리케이션이 백그라운드에서 시작되었습니다. 로그는 app.log 파일에서 확인할 수 있습니다."
    
    log "원격 서버 작업 완료"
EOF
log "AWS 서버 SSH 접속 및 작업 완료"

END_TIME=$(date +%s.%N)
EXECUTION_TIME=$(echo "$END_TIME - $START_TIME" | bc)
log "스크립트 실행 완료. 총 실행 시간: $(printf "%.3f" $EXECUTION_TIME) 초"

# 현재 디렉토리로 복귀
cd $CURRENT_DIR