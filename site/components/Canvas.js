import Image from "next/image";
import { useEffect, useRef, useState } from "react";
import useResizeObserver from "use-resize-observer";

const Canvas = ({ image, loadingImage }) => {
  const ref = useRef();
  const [canvasSize, setCanvasSize] = useState({ width: 1, height: 1 });
  const { width, height } = useResizeObserver({ ref });

  useEffect(() => {
    if (!!ref.current) {
      const rect = ref.current.children[0].getBoundingClientRect();
      setCanvasSize({ width: rect.width, height: rect.height });
    }
  }, [width, height, ref]);

  return (
    <>
      {image ? (
        <div className="overflow-hidden rounded-lg">
          <div className={"relative" + (loadingImage ? " blur-xl" : "")}>
            <div ref={ref}>
              <Image
                src={image.src}
                width={image.width}
                height={image.height}
                alt="Image being sonified"
                className="absolute top-0 left-0 rounded-md shadow-md"
              />
            </div>
            <canvas
              id="pixelsound"
              width={canvasSize.width}
              height={canvasSize.height}
              className="absolute top-0 left-0"
            ></canvas>
          </div>
        </div>
      ) : (
        <div className="animate-pulse bg-slate-100 w-48 h-48 rounded-md shadow-md"></div>
      )}
    </>
  );
};

export default Canvas;
