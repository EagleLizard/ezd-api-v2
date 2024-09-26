package ctrl

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/EagleLizard/jcd-api/gosrc/lib/config"
	"github.com/EagleLizard/jcd-api/gosrc/lib/logging"
	"github.com/go-chi/chi/v5"
)

func GetImage(
	cfg config.JcdApiConfigType,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()
		folderKey := chi.URLParam(r, "folderKey")
		imageKey := chi.URLParam(r, "imageKey")
		if cfg.JcdEnv != "DEV" {
			/*
				TODO: remove this block once real object store is
					hooked up
			*/
			logging.Logger.Fatal("GetImage dev error")
		}
		// fmt.Printf("%s / %s\n", folderKey, imageKey)
		sfsAddr := net.JoinHostPort(cfg.SfsHost, cfg.SfsPort)
		imgPath, err := url.JoinPath("img-v3", folderKey, imageKey)
		if err != nil {
			logging.Logger.Sugar().Error(err)
		}
		imgUrl := url.URL{
			Scheme: "http",
			Host:   sfsAddr,
			Path:   imgPath,
		}
		fmt.Printf("%s\n", imgUrl.String())
		resp, err := http.Get(imgUrl.String())
		if err != nil {
			logging.Logger.Sugar().Error(err)
		}
		fmt.Printf("%s\n", resp.Status)
		w.WriteHeader(http.StatusOK)
		/*
			TODO:
				- detect content-type
				- stream image to resp
		*/
	}
}
