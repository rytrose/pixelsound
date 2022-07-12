import { useRef, useCallback, useEffect } from "react";
import { WaveSurfer, WaveForm } from "wavesurfer-react";

const Waveform = ({ audio, loadingAudio }) => {
  const onWaveSurferReady = useCallback(() => {
    // TODO: determine if needed
  }, []);

  const onWaveSurferLoading = useCallback((data) => {
    // TODO: determine if needed
  }, []);

  const waveSurferRef = useRef();
  const onWaveSurferMount = useCallback(
    (waveSurfer) => {
      waveSurferRef.current = waveSurfer;
      if (waveSurferRef.current) {
        const ws = waveSurferRef.current;
        ws.on("ready", onWaveSurferReady);
        ws.on("loading", onWaveSurferLoading);
      }
    },
    [onWaveSurferReady, onWaveSurferLoading]
  );

  useEffect(() => {
    if (audio && waveSurferRef.current) {
      const ws = waveSurferRef.current;
      ws.loadBlob(audio);
    }
  }, [audio]);

  // Called from golang code when playing audio
  const updateWaveform = useCallback((progress) => {
    if (waveSurferRef.current) {
      const ws = waveSurferRef.current;
      ws.seekTo(progress);
    }
  }, []);

  useEffect(() => {
    if (!window.jsUpdateWaveform) window.jsUpdateWaveform = updateWaveform;
  }, [updateWaveform]);

  return (
    <div
      className={
        "h-24 rounded-md shadow-md  " +
        (!audio ? "w-48 animate-pulse bg-slate-100" : "w-full")
      }
    >
      <WaveSurfer onMount={onWaveSurferMount}>
        {!!audio && (
          <WaveForm
            id="waveform"
            cursorColor="transparent"
            height={96}
            barWidth={2}
            responsive={true}
            waveColor={"#64748b"}
            progressColor={"#000000"}
          ></WaveForm>
        )}
      </WaveSurfer>
    </div>
  );
};

export default Waveform;
