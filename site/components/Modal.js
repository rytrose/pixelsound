import { useEffect, useRef, useState } from "react";

const Modal = ({ onClose, children }) => {
  const dialogEl = useRef();
  const [dialogOpened, setDialogOpened] = useState(false);

  // Required to trigger the opacity transition :(
  // See question on SO: https://stackoverflow.com/q/72835953/3434708
  useEffect(() => {
    const dialog = dialogEl.current;
    if (!dialog.open && !dialogOpened) {
      dialog.showModal();
      dialog.classList.remove("opacity-0");
      dialog.classList.add("opacity-100");
      setDialogOpened(true);
    } else {
      dialog.classList.remove("opacity-0");
    }
  }, [dialogOpened]);

  return (
    <dialog
      ref={dialogEl}
      className={`transition duration-500 opacity-0 backdrop:bg-slate-300 rounded-xl`}
      onClose={onClose}
    >
      {children}
    </dialog>
  );
};

export default Modal;
