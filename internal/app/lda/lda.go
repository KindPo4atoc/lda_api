package lda

import (
	//"errors"
	"errors"
	"fmt"
	"math"

	"github.com/lda_api/internal/app/entity"
	"github.com/lda_api/internal/app/repository"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
)

// структура lda
type LDA struct {
	db            *repository.DataBase
	AccuracyModel float64     `json:"accuracy"`
	X             []float64   `json:"x"`
	Y             []float64   `json:"y"`
	Classes       []int       `json:"class"`
	LinearDisc    [][]float64 `json:"linear_disc"`
	ClassMeans    [][]float64 `json:"class_means"`
	CovInvMatrix  [][]float64
}

// конструктор
func New(db *repository.DataBase) *LDA {
	return &LDA{
		db: db,
	}
}

// метод для получения дисперсии
func GetDispersion(data []float64, mean float64) float64 {
	if data == nil {
		return 0.0
	}
	var sum float64 = 0.0
	for i := 0; i < len(data); i++ {
		var temp float64 = data[i] - mean
		sum = sum + math.Pow(temp, 2)
	}
	return sum / (float64(len(data) - 1))
}

// метод для получения ковариации
func GetCovariation(dataX []float64, meanX float64, dataY []float64, meanY float64) float64 {
	if dataX == nil && dataY == nil {
		logrus.Fatal("Data is Empty")
		return 0.0
	}
	var sum float64 = 0.0
	for i := 0; i < len(dataX); i++ {
		sum = sum + (dataX[i]-meanX)*(dataY[i]-meanY)
	}
	return sum / (float64(len(dataY) - 1))
}

// метод для получения ковариационной матрицы
func GetCovariationMatrix(mean []float64, dataCoef, dataRating []float64) [][]float64 {
	var covariationMatrix [][]float64
	covariation := GetCovariation(dataCoef, mean[1], dataRating, mean[0]) //проверить зависимость коэф
	for i := 0; i < len(mean); i++ {
		var temp []float64
		for j := 0; j < len(mean); j++ {
			if i == j {
				if i == 0 {
					temp = append(temp, GetDispersion(dataRating, mean[i]))
				} else {
					temp = append(temp, GetDispersion(dataCoef, mean[i]))
				}
			} else {
				temp = append(temp, covariation)
			}
		}
		covariationMatrix = append(covariationMatrix, temp)
	}
	return covariationMatrix
}

// метод для сложения матриц
func PlusMatrix(leftOperand, rightOperand [][]float64) [][]float64 {
	if leftOperand == nil {
		return rightOperand
	}
	for i := 0; i < len(leftOperand); i++ {
		for j := 0; j < len(leftOperand[i]); j++ {
			leftOperand[i][j] = leftOperand[i][j] + rightOperand[i][j]
		}
	}
	return leftOperand
}

// метод инициализации обучающих данных
func initData(data entity.ContextData) ([][]float64, []int) {
	var class = make([]int, len(data.Data))
	var dataForLearn [][]float64

	for i := 0; i < len(data.Data); i++ {
		var tmp []float64
		tmp = append(tmp, float64(data.Data[i].IncomeAnnum))
		tmp = append(tmp, float64(data.Data[i].LoanAmount))
		tmp = append(tmp, float64(data.Data[i].LoanTerm))
		tmp = append(tmp, float64(data.Data[i].CibilScore))

		dataForLearn = append(dataForLearn, tmp)
		if data.Data[i].LoanStatus == "Approved" {
			class[i] = 1
		} else {
			class[i] = 0
		}
	}
	fmt.Println(len(dataForLearn))
	fmt.Println(len(class))
	return dataForLearn, class
}

// метод возвращающий уникальные значения
func Unique(class []int) []int {
	var result []int
	for i := 0; i < len(class); i++ {
		temp := class[i]
		check := false
		for j := 0; j < len(result); j++ {
			if temp == result[j] {
				check = true
				break
			}
		}
		if !check {
			result = append([]int{temp}, result...)
		}
	}
	return result
}

// метод расчета средниз значений
func GetMeans(X [][]float64) []float64 {
	var result []float64

	for i := 0; i < len(X[0]); i++ {
		var temp = 0.0
		for j := 0; j < len(X); j++ {
			temp += X[j][i]
		}
		temp = temp / float64(len(X))
		result = append(result, temp)
	}
	return result
}

// метод возвращающий данные определенного класса
func GetClassValue(class int, classes []int, X [][]float64) [][]float64 {
	var result [][]float64
	for i := 0; i < len(classes); i++ {
		if classes[i] == class {

			result = append(result, X[i])
		}
	}
	return result
}

// метод умножения матриц
func ProdsMatrix(a, b [][]float64) ([][]float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return nil, fmt.Errorf("одна из матриц пустая")
	}

	// Получаем размеры матриц
	aRows := len(a)
	aCols := len(a[0])
	bRows := len(b)
	bCols := len(b[0])

	// Проверяем совместимость размеров
	if aCols != bRows {
		return nil, fmt.Errorf("несовместимые размеры матриц: a[%d×%d] * b[%d×%d]",
			aRows, aCols, bRows, bCols)
	}

	// Создаем результирующую матрицу
	result := make([][]float64, aRows)
	for i := range result {
		result[i] = make([]float64, bCols)
	}

	// Выполняем умножение
	for i := 0; i < aRows; i++ {
		for j := 0; j < bCols; j++ {
			for k := 0; k < aCols; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}

	return result, nil
}

// метод умноения матрицы на число
func ProdMatrixWithValue(a float64, b [][]float64) [][]float64 {
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[i]); j++ {
			b[i][j] = b[i][j] * a
		}
	}
	return b
}

// вычитание матриц
func MinusMatrix(leftOperand, rightOperand []float64) []float64 {
	var result []float64
	for i := 0; i < len(leftOperand); i++ {
		temp := leftOperand[i] - rightOperand[i]
		result = append(result, temp)
	}
	return result
}

// метод транспонирования матрицы
func T(matrix [][]float64) [][]float64 {
	rows := len(matrix)
	if rows == 0 {
		return [][]float64{}
	}
	cols := len(matrix[0])

	// Создаем новую матрицу с обратными размерами
	transposed := make([][]float64, cols)
	for i := range transposed {
		transposed[i] = make([]float64, rows)
	}

	// Заполняем транспонированную матрицу
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			transposed[i][j] = matrix[j][i]
		}
	}

	return transposed
}

// метод получения обратной матрицы
func inverseMatrix(matrix [][]float64) ([][]float64, error) {
	n := len(matrix)
	if n == 0 {
		return nil, fmt.Errorf("матрица пустая")
	}

	// Создаем расширенную матрицу [A|I]
	augmented := make([][]float64, n)
	for i := range augmented {
		augmented[i] = make([]float64, 2*n)
		copy(augmented[i][:n], matrix[i])
		augmented[i][n+i] = 1
	}

	// Прямой ход метода Гаусса
	for col := 0; col < n; col++ {
		// Частичный выбор ведущего элемента
		pivot := col
		maxVal := math.Abs(augmented[col][col])
		for i := col + 1; i < n; i++ {
			if abs := math.Abs(augmented[i][col]); abs > maxVal {
				maxVal = abs
				pivot = i
			}
		}

		if maxVal < 1e-12 {
			return nil, fmt.Errorf("матрица вырожденная, обратной не существует")
		}

		// Перестановка строк
		if pivot != col {
			augmented[col], augmented[pivot] = augmented[pivot], augmented[col]
		}

		// Нормализация текущей строки
		divisor := augmented[col][col]
		for j := col; j < 2*n; j++ {
			augmented[col][j] /= divisor
		}

		// Обнуление элементов в текущем столбце
		for i := 0; i < n; i++ {
			if i != col {
				factor := augmented[i][col]
				for j := col; j < 2*n; j++ {
					augmented[i][j] -= factor * augmented[col][j]
				}
			}
		}
	}

	// Извлекаем обратную матрицу из правой части
	inverse := make([][]float64, n)
	for i := range inverse {
		inverse[i] = make([]float64, n)
		copy(inverse[i], augmented[i][n:])
	}

	return inverse, nil
}

// метод сортировки собственных значений
func SortDesc(values []float64) ([]float64, []int) {
	//var resultArr []float64
	copyArr := make([]float64, len(values))
	for i := range values {
		copyArr[i] = values[i]
	}
	var idxs []int
	for i := 0; i < len(values); i++ {
		var lastValue float64
		if i != len(values)-1 {
			for j := 0; j+1 < len(values)-i; j++ {
				if values[j] < values[j+1] {
					temp := values[j]
					values[j] = values[j+1]
					values[j+1] = temp
				}
			}
		}
		lastValue = values[len(values)-1-i]
		fmt.Println(lastValue)
		for j := 0; j < len(copyArr); j++ {
			if lastValue == copyArr[j] {
				idxs = append([]int{j}, idxs...)
				//idxs = append(idxs, j)
				break
			}
		}
	}
	return values, idxs
}

// метод сортировки собственных векторов
func SortVectors(vectors [][]float64, idxs []int) [][]float64 {
	result := make([][]float64, len(vectors))
	for i := 0; i < len(idxs); i++ {
		result[i] = vectors[idxs[i]]
	}
	return result
}

// метод для получения коэффициентов
func GetCoefficients(data [][]float64, n int) ([][]float64, error) {
	if len(data) < n {
		return nil, fmt.Errorf("Компонентов больше чем размерность векторов")
	}

	if n == len(data) {
		return data, nil
	} else {
		result := make([][]float64, n)
		for i := 0; i < n; i++ {
			result[i] = data[i]
		}
		return result, nil
	}
}

// скалярное произведение векторов
func SampleProdMatrix(leftOperand, rightOperand [][]float64) ([][]float64, error) {
	if len(leftOperand) != len(rightOperand) || len(leftOperand[0]) != len(rightOperand[0]) {
		return nil, fmt.Errorf("Матрицы имеют разные размерности")
	}
	leftOperand[0][0] = leftOperand[0][0] * rightOperand[0][0]
	leftOperand[0][1] = leftOperand[0][1] * rightOperand[0][1]
	return leftOperand, nil
}

// метод расчета данных для пространства lda
func (lda *LDA) TransformData(data [][]float64, init bool) ([][]float64, error) {
	result, err := ProdsMatrix(data, T(lda.LinearDisc))
	if err != nil {
		return nil, err
	}
	if init {
		for i := 0; i < len(result); i++ {
			lda.X = append(lda.X, result[i][0])
			lda.Y = append(lda.Y, result[i][1])
		}
	}
	return result, nil
}

// метод для осуществления прогноза
func (lda *LDA) Predict(data [][]float64) (int, [][]float64, [][]float64, error) {
	ldaData, err := lda.TransformData(data, false)

	if err != nil {
		return 0.0, nil, nil, err
	}
	var distance [][]float64
	for i := 0; i < len(lda.ClassMeans); i++ {
		var diff [][]float64
		for j := 0; j < len(ldaData); j++ {
			diff = append(diff, MinusMatrix(ldaData[j], lda.ClassMeans[i]))
		}
		var mahalanobisDistance []float64
		prodDiffCov, err := ProdsMatrix(diff, lda.CovInvMatrix)
		if err != nil {
			return 0.0, nil, nil, err
		}
		resultProd, err := SampleProdMatrix(prodDiffCov, diff)
		if err != nil {
			return 0.0, nil, nil, err
		}
		tmp := 0.0
		for j := 0; j < len(resultProd[0]); j++ {
			tmp += resultProd[0][j]
		}
		mahalanobisDistance = append(mahalanobisDistance, tmp)
		distance = append(distance, mahalanobisDistance)
	}
	distanceT := T(distance)

	min := math.MaxFloat64
	indxMin := -1
	for i := 0; i < len(distanceT[0]); i++ {
		if min > distanceT[0][i] {
			min = distanceT[0][i]
			indxMin = i
		}
	}

	return indxMin, distance, ldaData, nil
}

// метод для получения точности модели
func (lda *LDA) GetAccuracyModel(data [][]float64, y []int) error {
	count := 0.0
	for i := 0; i < len(data); i++ {
		var tmp [][]float64
		tmp = append(tmp, data[i])
		predict, _, _, err := lda.Predict(tmp)
		if err != nil {
			return err
		}
		//fmt.Println("distance: ", d)
		//fmt.Printf("predict - %d ----- class - %d\n", predict, y[i])
		if predict == y[i] {
			count += 1
		}
	}
	lda.AccuracyModel = (count / float64(len(data))) * 100.0
	return nil
}
func (lda *LDA) GetStringCoef() string {
	resultStr := ""
	for i := 0; i < len(lda.LinearDisc); i++ {
		for j := 0; j < len(lda.LinearDisc[i]); j++ {
			if j == 0 {
				resultStr += "| " + fmt.Sprint(lda.LinearDisc[i][j]) + "\t"
			}
			if j == len(lda.LinearDisc[i])-1 {
				resultStr += fmt.Sprint(lda.LinearDisc[i][j]) + " |\n"
			} else {
				resultStr += fmt.Sprint(lda.LinearDisc[i][j]) + "\t"
			}
		}
	}
	return resultStr
}

// главный метод для обучения модели
func (lda *LDA) FitModel() error {
	dataContext, err := lda.db.Data().SelectAllLearnData()
	if err != nil {
		return err
	}
	data, class := initData(dataContext)
	lda.Classes = class
	nFeatures := len(data[0])
	n_components := 2
	classLabel := Unique(class)
	nClasses := len(classLabel)

	meanAllData := GetMeans(data)
	fmt.Println("classlabel")
	fmt.Println(classLabel)
	fmt.Println(nFeatures)
	fmt.Println(nClasses)
	fmt.Println(meanAllData)

	var S_W [][]float64
	var S_B [][]float64

	var classMeans [][]float64
	for i := 0; i < len(classLabel); i++ {
		dataClass := GetClassValue(classLabel[i], class, data)
		dataMeanClass := GetMeans(dataClass)
		classMeans = append(classMeans, dataMeanClass)

		var S [][]float64
		for j := 0; j < len(dataClass); j++ {
			diff := MinusMatrix(dataClass[j], dataMeanClass)
			S = append(S, diff)
		}
		sT := T(S)
		S, err := ProdsMatrix(sT, S)
		if err != nil {
			return err
		}
		S_W = PlusMatrix(S_W, S)

		countData := len(dataClass)
		tmpDiff := MinusMatrix(dataMeanClass, meanAllData)
		var tmp [][]float64
		tmp = append(tmp, tmpDiff)
		meanDiff := T(tmp)
		S, err = ProdsMatrix(meanDiff, tmp)
		if err != nil {
			return err
		}
		S_B = PlusMatrix(S_B, ProdMatrixWithValue(float64(countData), S))

	}
	inv, err := inverseMatrix(S_W)

	if err != nil {
		return err
	}

	A, err := ProdsMatrix(inv, S_B)
	if err != nil {
		return err
	}
	var tmpA []float64

	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A[i]); j++ {
			tmpA = append(tmpA, A[i][j])
		}
	}
	B := mat.NewDense(4, 4, tmpA)
	var eig mat.Eigen
	if !eig.Factorize(B, mat.EigenRight) {
		err := errors.New("error in calculating eigenvalues")
		return err
	}
	var eigValues []float64
	values := eig.Values(nil)
	if values == nil {
		err := errors.New("couldn't get eigenvalues")
		return err
	}
	for i := 0; i < 4; i++ {
		eigValues = append(eigValues, float64(real(values[i])))
	}

	// Получаем собственные векторы
	var vectors mat.CDense
	eig.VectorsTo(&vectors)
	var eigVectors [][]float64

	for i := 0; i < 4; i++ {
		var temp []float64
		for j := 0; j < 4; j++ {
			temp = append(temp, float64(real(vectors.At(i, j))))
		}
		eigVectors = append(eigVectors, temp)
	}

	eigVectors = T(eigVectors)
	_, idxs := SortDesc(eigValues)
	sortEigVectors := SortVectors(eigVectors, idxs)
	lda.LinearDisc, err = GetCoefficients(sortEigVectors, n_components)
	if err != nil {
		return err
	}
	// преобразование средних значений для пространства lda
	lda.ClassMeans, err = ProdsMatrix(classMeans, T(lda.LinearDisc))
	if err != nil {
		return err
	}

	convLdaData, err := lda.TransformData(data, true)
	if err != nil {
		return err
	}
	convLdaDataT := T(convLdaData)
	ldaMeanVal := GetMeans(convLdaData)
	covariationLdaMatrix := GetCovariationMatrix(ldaMeanVal, convLdaDataT[1], convLdaDataT[0])
	fmt.Println(lda.LinearDisc)
	lda.CovInvMatrix, err = inverseMatrix(covariationLdaMatrix)
	if err != nil {
		return err
	}
	test, err := lda.db.Data().SelectTestData()
	if err != nil {
		return err
	}
	dataTest, yTest := initData(test)
	err = lda.GetAccuracyModel(dataTest, yTest)
	logrus.Info("------Параметры модели------")
	logrus.Info("Точность модели ", lda.AccuracyModel)
	logrus.Info("Коэффициенты модели:")
	fmt.Println(lda.GetStringCoef())
	if err != nil {
		return err
	}
	return nil
}
