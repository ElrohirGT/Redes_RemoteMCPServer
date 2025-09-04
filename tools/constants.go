package tools

import (
	"log"
	"os"
	"time"
)

var l = log.New(os.Stderr, "", log.Lshortfile|log.Ldate|log.Ltime)

const API_BASE_URL = "https://api.magicthegathering.io"
const CACHE_DURATION = 12 * time.Hour
