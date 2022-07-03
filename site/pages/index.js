import Head from "next/head";
import Script from "next/script";
import { useEffect, useState } from "react";
import Button from "../components/Button";
import Canvas from "../components/Canvas";
import Controls from "../components/Controls";
import Header from "../components/Header";
import Loading from "../components/Loading";
import Modal from "../components/Modal";

const Home = () => {
  const [loading, setLoading] = useState(true);
  const [started, setStarted] = useState(false);
  const [image, setImage] = useState();
  const [loadingImage, setLoadingImage] = useState(true);

  const start = () => {
    setStarted(true);
    // Made available globally by golang code
    window.golangSetup();
  };

  // Called when golang code has finished populating the window
  const golangReady = () => {
    setLoading(false);
  };

  // Called when golang code has finished updating the image
  const imageUpdated = () => {
    setLoadingImage(false);
  };

  // Called when golang code has finished updating the audio
  const audioUpdated = () => {
    // TODO
    console.log("audio updated");
  };

  // Setup functions exposed to golang on window
  useEffect(() => {
    window.jsGolangReady = golangReady;
    window.jsImageUpdated = imageUpdated;
    window.jsAudioUpdated = audioUpdated;
  }, []);

  const onImageChange = (e) => {
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
  };

  const onAudioChange = (e) => {
    const input = e.target;
    if (input.files && input.files[0]) {
      // Mark the audio as loading until golang has updated

      // Read the file and update in JS and golang
      const reader = new FileReader();
      reader.onload = (e) => {
        // Display the audio

        // Made available globally by golang code
        window.golangUpdateAudio(e.target.result);
      };
      reader.readAsDataURL(input.files[0]);
    }
  };
  const onModeChange = (e) => console.log(e.target);

  return (
    <>
      <Script src="https://code.jquery.com/jquery-3.6.0.slim.min.js"></Script>
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
        <Header className="border-b-[1px] border-slate-300"></Header>
        <Canvas
          loadingImage={loadingImage}
          image={image}
          className="border-b-[1px] border-slate-300"
        ></Canvas>
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
