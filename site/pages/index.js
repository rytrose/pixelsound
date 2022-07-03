import Head from "next/head";
import Script from "next/script";
import { useState } from "react";
import Button from "../components/Button";
import Canvas from "../components/Canvas";
import Controls from "../components/Controls";
import Header from "../components/Header";
import Modal from "../components/Modal";

const Home = () => {
  const [started, setStarted] = useState(false);
  const [image, setImage] = useState();

  const start = () => {
    setStarted(true);
    // Made available globally by golang code
    window.golangRun();
  };

  const onImageChange = (e) => {
    const input = e.target;
    if (input.files && input.files[0]) {
      const reader = new FileReader();
      reader.onload = (e) => {
        const img = new Image();
        img.onload = () => setImage(img);
        img.src = e.target.result;
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
        <div className="text-center">
          <p className="py-2 font-serif text-4xl text-black">pixelsound</p>
          <p className="text-sm text-stone-600">
            An image sonification playground
          </p>
          <form className="pt-4" method="dialog">
            <Button>Start</Button>
          </form>
        </div>
      </Modal>

      <div className={started ? "visible" : "invisible"}>
        <Header className="border-b-[1px] border-slate-300"></Header>
        <Canvas
          onImageChange={onImageChange}
          image={image}
          className="border-b-[1px] border-slate-300"
        ></Canvas>
        <Controls onModeChange={onModeChange}></Controls>
      </div>
    </>
  );
};

export default Home;
