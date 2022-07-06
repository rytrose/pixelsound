import dynamic from "next/dynamic";
import Head from "next/head";
import Script from "next/script";
import { useEffect, useState, useCallback } from "react";
import Button from "../components/Button";
import Canvas from "../components/Canvas";
import Controls from "../components/Controls";
import Header from "../components/Header";
import Loading from "../components/Loading";
import Modal from "../components/Modal";

const Waveform = dynamic(() => import("../components/Waveform"), {
  ssr: false,
});

const Home = () => {
  const [loading, setLoading] = useState(true);
  const [started, setStarted] = useState(false);
  const [image, setImage] = useState();
  const [audio, setAudio] = useState();
  const [loadingImage, setLoadingImage] = useState(true);
  const [loadingAudio, setLoadingAudio] = useState(true);

  const start = useCallback(() => {
    // Made available globally by golang code
    window.golangSetup();
  }, []);

  // Called when golang code has finished populating the window
  const golangReady = useCallback(() => {
    setLoading(false);
  }, []);

  // Called when golang code has finished setting up
  const golangSetup = useCallback(() => {
    setStarted(true);
  }, []);

  // Called when golang code has finished updating the image
  const imageUpdated = useCallback(() => {
    setLoadingImage(false);
  }, []);

  // Called when golang code has finished updating the audio
  const audioUpdated = useCallback(() => {
    setLoadingAudio(false);
  }, []);

  // Setup functions exposed to golang on window
  useEffect(() => {
    if (!window.jsGolangReady) window.jsGolangReady = golangReady;
    if (!window.jsGolangSetup) window.jsGolangSetup = golangSetup;
    if (!window.jsImageUpdated) window.jsImageUpdated = imageUpdated;
    if (!window.jsAudioUpdated) window.jsAudioUpdated = audioUpdated;
    // Run golang logic that depends on elements
    if (!loading && started) {
      // Made available globally by golang code
      window.golangRun();
    }
  }, [golangReady, golangSetup, imageUpdated, audioUpdated, loading, started]);

  const onImageChange = useCallback((e) => {
    const input = e.target;
    if (input.files && input.files[0]) {
      // Reset the image to placeholder
      setImage(undefined);
      // Mark the image as loading until golang has updated
      setLoadingImage(true);

      // Read the file and update in JS and golang
      const reader = new FileReader();
      reader.onload = (e) => {
        // Load the image to read dimensions
        const img = new Image();
        img.onload = () => setImage(img);
        img.src = e.target.result;

        // Made available globally by golang code
        window.golangUpdateImage(e.target.result);
      };
      reader.readAsDataURL(input.files[0]);
    }
  }, []);

  const onAudioChange = useCallback((e) => {
    const input = e.target;
    if (input.files && input.files[0]) {
      // Mark the audio as loading until golang has updated
      setLoadingAudio(true);

      // Display the audio
      setAudio(input.files[0]);

      // Read the file and update in JS and golang
      const reader = new FileReader();
      reader.onload = (e) => {
        // Made available globally by golang code
        window.golangUpdateAudio(e.target.result);
      };
      // Read the audio for golang
      reader.readAsDataURL(input.files[0]);
    }
  }, []);

  const onModeChange = useCallback((e) => console.log(e.target), []);

  return (
    <>
      <Script src="/pixelsound.js"></Script>

      <Head>
        <title>pixelsound</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Modal onClose={start}>
        <div className="flex flex-col justify-center items-center">
          <p className="py-2 font-serif text-4xl text-black">pixelsound</p>
          <p className="text-sm text-stone-600">
            An image sonification playground
          </p>
          <div className="flex items-center h-16">
            {loading ? (
              <Loading className="w-16" />
            ) : (
              <form className="" method="dialog">
                <Button>Start</Button>
              </form>
            )}
          </div>
        </div>
      </Modal>

      <div className={started ? "visible" : "invisible"}>
        <div className="border-b-[1px] border-slate-300">
          <Header />
        </div>
        <div className="p-3 border-b-[1px] border-slate-300">
          <div className="flex flex-col gap-3 max-w-lg items-center mx-auto">
            <Canvas loadingImage={loadingImage} image={image}></Canvas>
            <Waveform loadingAudio={loadingAudio} audio={audio}></Waveform>
          </div>
        </div>
        <Controls
          onImageChange={onImageChange}
          onAudioChange={onAudioChange}
          onModeChange={onModeChange}
        ></Controls>
      </div>
    </>
  );
};

export default Home;
