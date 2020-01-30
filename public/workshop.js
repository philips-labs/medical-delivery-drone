const displaySize = { width: 640, height: 480 }
let labeledFaceDescriptors = [];

const init = async () => {
  const MODEL_URL = '/models';

  await faceapi.loadSsdMobilenetv1Model(MODEL_URL)
  await faceapi.loadFaceLandmarkModel(MODEL_URL)
  await faceapi.loadFaceRecognitionModel(MODEL_URL)

  const canvas = document.getElementById('overlay')
  const labels = ['sjaak', 'richard', 'arvid', 'mara', 'frank', 'wilma', 'marcin', 'deha', 'slave', 'louis', 'jeroen']

  labeledFaceDescriptors = await Promise.all(
    labels.map(async label => {
      // fetch image data from urls and convert blob to HTMLImage element
      const imgUrl = `labels/${label}.png`
      const img = await faceapi.fetchImage(imgUrl)

      // detect the face with the highest score in the image and compute it's landmarks and face descriptor
      const fullFaceDescription = await faceapi.detectSingleFace(img).withFaceLandmarks().withFaceDescriptor()

      if (!fullFaceDescription) {
        throw new Error(`no faces detected for ${label}`)
      }

      const faceDescriptors = [fullFaceDescription.descriptor]
      return new faceapi.LabeledFaceDescriptors(label, faceDescriptors)
    })
  );
  
  faceapi.matchDimensions(canvas, displaySize)

  const videoEl = document.getElementById('inputVideo')
  navigator.getUserMedia(
    { video: {} },
    stream => videoEl.srcObject = stream,
    err => console.error(err)
  )
}

async function onPlay(videoEl) {
  const input = document.getElementById('inputVideo')
  const canvas = document.getElementById('overlay')

  let fullFaceDescriptions = await faceapi.detectAllFaces(input).withFaceLandmarks().withFaceDescriptors()
  if (fullFaceDescriptions.length > 0) {
    await detectPersons(fullFaceDescriptions, canvas);
    fullFaceDescriptions = faceapi.resizeResults(fullFaceDescriptions, displaySize);
  }

  setTimeout(() => onPlay(videoEl))
}

async function detectPersons(fullFaceDescriptions, canvas) {
  const context = canvas.getContext('2d');
  const maxDescriptorDistance = 0.6
  const faceMatcher = new faceapi.FaceMatcher(labeledFaceDescriptors, maxDescriptorDistance)
  const results = fullFaceDescriptions.map(fd => faceMatcher.findBestMatch(fd.descriptor))

  context.clearRect(0, 0, canvas.width, canvas.height);
  results.forEach((bestMatch, i) => {
    const box = fullFaceDescriptions[i].detection.box
    const text = bestMatch.toString()
    const drawBox = new faceapi.draw.DrawBox(box, { label: text })
    drawBox.draw(canvas)
  })
}

init().then(() => console.log("Done"));
